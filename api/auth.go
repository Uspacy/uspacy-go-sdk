package api

import (
	"encoding/json"
	"github.com/Uspacy/uspacy-go-sdk/auth"
	"log"
)

func (us *Uspacy) getToken() string {
	return us.bearerToken

}

func (us *Uspacy) refreshToken() {
	var refresh auth.RefreshOutput

	body, err := us.doPostEmptyHeaders(us.mainHost+auth.VersionUrl+auth.RefreshTokenUrl, nil)
	if err != nil {
		log.Fatal("error while trying to refresh token: ", err)
	}
	err = json.Unmarshal(body, &refresh)
	if err != nil {
		log.Fatal("error while trying to parse token: ", err)
	}
	us.bearerToken = refresh.Jwt

}
