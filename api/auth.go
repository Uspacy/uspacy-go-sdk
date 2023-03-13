package api

import (
	"encoding/base64"
	"encoding/json"
	"github.com/Uspacy/uspacy-go-sdk/auth"
	"log"
	"strings"
	"time"
)

func (us *Uspacy) getToken() string {
	if time.Now().Unix() > us.unixExpTime {
		us.isExpired = true
		token, expDate := us.refreshToken()
		us.bearerToken = token
		us.unixExpTime = time.Now().Local().Add(time.Second * time.Duration(expDate)).Unix()
		us.isExpired = false
	}

	return us.bearerToken

}

func (us *Uspacy) refreshToken() (string, int64) {
	var refresh auth.RefreshOutput

	body, err := us.doPostEmptyHeaders(mainHost+auth.VersionUrl+auth.RefreshTokenUrl, nil)
	if err != nil {
		log.Fatal("error while trying to refresh token:", err)
	}
	err = json.Unmarshal(body, &refresh)
	if err != nil {
		log.Fatal("error while trying to parse token: ", err)
	}

	return refresh.Jwt, refresh.ExpireInSeconds

}

func setExp(token string) int64 {
	jwtData, err := getJWTData(token)
	if err != nil {
		log.Fatal("token is invalid")
	}
	return jwtData.Exp

}

func getJWTData(token string) (tokenData auth.JwtClaims, err error) {
	prefix := "Bearer "
	reqToken := strings.TrimPrefix(token, prefix)
	_string := strings.Split(strings.Join(strings.Split(reqToken, " "), "."), ".")
	for i := 0; i < len(_string); i++ {
		decoded, _ := base64.RawURLEncoding.DecodeString(_string[i])
		err = json.Unmarshal(decoded, &tokenData)
		if err != nil {
			return tokenData, err
		}

	}
	return
}
