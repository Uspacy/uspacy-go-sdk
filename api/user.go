package api

import (
	"encoding/json"
	"fmt"

	"github.com/Uspacy/uspacy-go-sdk/user"
)

// CreateActiveUser returns created users
func (us *Uspacy) CreateActiveUsers(entity []user.User) (object []user.User, err error) {
	body, err := us.doPostEmptyHeaders(us.buildURL(user.VersionUrl, user.CreateActiveUser), entity)
	if err != nil {
		return object, err
	}
	return object, json.Unmarshal(body, &object)
}

// PatchUsers patch user by Id and return it
func (us *Uspacy) PatchUsers(entity user.User, userId string) (object user.User, err error) {
	body, err := us.doPatchEmptyHeaders(us.buildURL(user.VersionUrl, fmt.Sprintf(user.UserUrl, userId)), entity)
	if err != nil {
		return object, err
	}
	return object, json.Unmarshal(body, &object)
}
