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

func (us *Uspacy) getToken() string {
	return us.bearerToken
}

func (us *Uspacy) refreshToken() {
	var refresh auth.RefreshOutput
	body, err := us.doRaw(
		fmt.Sprintf("%s%s/%s/%s", "https://", us.UnmarshalTokenData().Domain, auth.VersionUrl, auth.RefreshTokenUrl),
		http.MethodPost,
		emptyHeaders,
		nil)
	if err != nil {
		log.Fatal("error while trying to refresh token: ", err)
	}
	err = json.Unmarshal(body, &refresh)
	if err != nil {
		log.Fatal("error while trying to parse token: ", err)
	}
	us.bearerToken = refresh.Jwt
}

func (us *Uspacy) UnmarshalTokenData() (tokenData auth.TokenData) {
	strings := strings.Split(strings.Join(strings.Split(us.bearerToken, " "), "."), ".")
	for _, _string := range strings {
		decoded, _ := base64.RawURLEncoding.DecodeString(_string)
		json.Unmarshal(decoded, &tokenData)
	}
	return
}
