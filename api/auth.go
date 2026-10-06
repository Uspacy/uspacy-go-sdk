package api

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Uspacy/uspacy-go-sdk/v2/auth"
)

// tokenRefresh refreshes the bearer token, carrying ctx on the refresh request (and its
// retry waits) so a done ctx aborts the refresh instead of running to completion.
// Concurrent callers are deduplicated: only one refresh request is executed and its
// result is shared.
func (us *Uspacy) tokenRefresh(ctx context.Context) (string, error) {
	res, err, _ := us.refreshFlight.Do("refresh", func() (any, error) {
		return us.refreshTokenOnce(ctx)
	})
	if err != nil {
		return "", err
	}
	return res.(string), nil
}

func (us *Uspacy) refreshTokenOnce(ctx context.Context) (string, error) {
	var refresh auth.RefreshOutput
	jwt, err := us.UnmarshalTokenData()
	if err != nil {
		return "", err
	}
	body, _, err := us.doRequest(
		ctx,
		fmt.Sprintf("%s%s/%s/%s", "https://", jwt.Domain, auth.VersionUrl, auth.RefreshTokenUrl),
		http.MethodPost,
		headersMap,
		nil,
		true) // skip token refresh: never refresh the token to refresh the token
	if err != nil {
		return "", err
	}
	if err := json.Unmarshal(body, &refresh); err != nil {
		return "", err
	}
	us.mu.Lock()
	us.bearerToken = refresh.Jwt
	us.refreshToken = refresh.RefreshToken
	us.mu.Unlock()
	return refresh.Jwt, nil
}

// TokenRefresh refreshes the bearer token, carrying ctx on the refresh request (and its
// retry waits) so a done ctx aborts the refresh instead of running to completion.
func (us *Uspacy) TokenRefresh(ctx context.Context) (string, error) {
	return us.tokenRefresh(ctx)
}

func (us *Uspacy) UnmarshalTokenData() (tokenData auth.JwtClaims, err error) {
	us.mu.RLock()
	token := us.bearerToken
	us.mu.RUnlock()

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
