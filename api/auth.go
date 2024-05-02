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
	body, _, err := us.doRaw(
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
	us.bearerToken = refresh.Jwt
	us.RefreshToken = refresh.RefreshToken
	return refresh.Jwt, nil
}

func (us *Uspacy) UnmarshalTokenData() (tokenData auth.JwtClaims, err error) {
	strings := strings.Split(strings.Join(strings.Split(us.bearerToken, " "), "."), ".")
	for _, _string := range strings {
		decoded, _ := base64.RawURLEncoding.DecodeString(_string)
		json.Unmarshal(decoded, &tokenData)
	}
	return
}
