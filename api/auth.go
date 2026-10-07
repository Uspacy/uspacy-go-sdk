package api

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Uspacy/uspacy-go-sdk/v2/auth"
)

// refreshTimeout bounds a shared token refresh, which runs detached from any single
// caller's ctx.
const refreshTimeout = time.Minute

// tokenRefresh refreshes the bearer token unless it already differs from staleToken, the
// token that got the 401 (another caller refreshed it meanwhile); an empty staleToken
// forces a refresh. Concurrent callers share one refresh request. That request does not
// inherit cancellation from whichever caller started it, so one caller giving up does not
// fail the refresh for the others; each caller still stops waiting as soon as its own ctx
// is done.
func (us *Uspacy) tokenRefresh(ctx context.Context, staleToken string) (string, error) {
	if token := us.currentToken(); staleToken != "" && token != staleToken {
		return token, nil
	}
	ch := us.refreshFlight.DoChan("refresh", func() (any, error) {
		refreshCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), refreshTimeout)
		defer cancel()
		return us.refreshTokenOnce(refreshCtx, staleToken)
	})
	select {
	case res := <-ch:
		if res.Err != nil {
			return "", res.Err
		}
		return res.Val.(string), nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

func (us *Uspacy) refreshTokenOnce(ctx context.Context, staleToken string) (string, error) {
	// Re-check inside the flight: a refresh that finished between the caller's check and
	// this one has already replaced the stale token.
	if token := us.currentToken(); staleToken != "" && token != staleToken {
		return token, nil
	}
	var refresh auth.RefreshOutput
	jwt, err := us.UnmarshalTokenData()
	if err != nil {
		return "", err
	}
	// The refresh endpoint authenticates with the refresh token, not the access token,
	// which may already have expired by the time a 401 triggers this refresh.
	_, refreshToken := us.Tokens()
	headers := mergeHeaders(nil, map[string]string{"Authorization": tokenPrefix + refreshToken})
	body, _, err := us.doRequest(
		ctx,
		fmt.Sprintf("%s%s/%s/%s", "https://", jwt.Domain, auth.VersionUrl, auth.RefreshTokenUrl),
		http.MethodPost,
		headers,
		nil,
		true) // skip token refresh: never refresh the token to refresh the token
	if err != nil {
		return "", err
	}
	if err := json.Unmarshal(body, &refresh); err != nil {
		return "", err
	}
	if refresh.Jwt == "" {
		return "", fmt.Errorf("token refresh: response has no jwt")
	}
	us.mu.Lock()
	us.bearerToken = refresh.Jwt
	// Keep the current refresh token if the response does not rotate it.
	if refresh.RefreshToken != "" {
		us.refreshToken = refresh.RefreshToken
	}
	us.mu.Unlock()
	return refresh.Jwt, nil
}

// Tokens returns the current access and refresh tokens as one consistent pair, so a
// caller persisting them never mixes tokens from two different refreshes. Safe for
// concurrent use.
func (us *Uspacy) Tokens() (access, refresh string) {
	us.mu.RLock()
	defer us.mu.RUnlock()
	return us.bearerToken, us.refreshToken
}

// TokenRefresh forces a bearer token refresh and returns the new token. It returns as
// soon as ctx is done; a refresh already in flight is shared with concurrent callers.
func (us *Uspacy) TokenRefresh(ctx context.Context) (string, error) {
	return us.tokenRefresh(ctx, "")
}

// UnmarshalTokenData decodes the claims of the current access token. It does not verify
// the token's signature.
func (us *Uspacy) UnmarshalTokenData() (tokenData auth.JwtClaims, err error) {
	token := us.currentToken()

	// JWT format: header.payload.signature
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return tokenData, fmt.Errorf("invalid JWT token format: expected 3 parts, got %d", len(parts))
	}

	// Decode payload (second part)
	decoded, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return tokenData, fmt.Errorf("failed to decode JWT payload: %w", err)
	}

	err = json.Unmarshal(decoded, &tokenData)
	if err != nil {
		return tokenData, fmt.Errorf("failed to unmarshal JWT claims: %w", err)
	}

	return tokenData, nil
}
