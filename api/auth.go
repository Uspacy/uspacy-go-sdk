package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Uspacy/uspacy-go-sdk/auth"
)

func (us *Uspacy) TokenRefresh() (string, error) {
	var refresh auth.RefreshOutput
	jwt, err := us.UnmarshalTokenData()
	if err != nil {
		return "", err
	}
	body, _, err := us.doRawSkipRefresh(
		fmt.Sprintf("%s%s/%s/%s", "https://", jwt.Domain, auth.VersionUrl, auth.RefreshTokenUrl),
		http.MethodPost,
		headersMap,
		nil)
	if err != nil {
		return "", err
	}
	err = json.Unmarshal(body, &refresh)
	if err != nil {
		return "", err
	}
	us.mu.Lock()
	us.bearerToken = refresh.Jwt
	us.RefreshToken = refresh.RefreshToken
	us.mu.Unlock()
	return refresh.Jwt, nil
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
