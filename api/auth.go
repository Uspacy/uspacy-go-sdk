package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/Uspacy/uspacy-go-sdk/auth"
)

func (us *Uspacy) refreshToken() error {
	var refresh auth.RefreshOutput
	jwt, err := us.UnmarshalTokenData()
	if err != nil {
		log.Fatal("error while trying to unmarshal token: ", err)
		return err
	}
	body, err := us.doRaw(
		fmt.Sprintf("%s%s/%s/%s", "https://", jwt.Domain, auth.VersionUrl, auth.RefreshTokenUrl),
		http.MethodPost,
		headersMap,
		nil)
	if err != nil {
		log.Fatal("error while trying to refresh token: ", err)
		return err
	}
	err = json.Unmarshal(body, &refresh)
	if err != nil {
		log.Fatal("error while trying to parse token: ", err)
		return err
	}
	us.bearerToken = refresh.Jwt
	return nil
}

func (us *Uspacy) UnmarshalTokenData() (tokenData auth.JwtClaims, err error) {
	strings := strings.Split(strings.Join(strings.Split(us.bearerToken, " "), "."), ".")
	for _, _string := range strings {
		decoded, _ := base64.RawURLEncoding.DecodeString(_string)
		err = json.Unmarshal(decoded, &tokenData)
	}
	return
}
