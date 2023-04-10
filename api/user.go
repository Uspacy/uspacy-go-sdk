package api

import (
	"encoding/json"
	"fmt"

	"github.com/Uspacy/uspacy-go-sdk/user"
)

// GetAllUsers gets all users
func (us *Uspacy) GetAllUsers() (object []user.User, err error) {
	body, err := us.doGetEmptyHeaders(us.buildURL(user.VersionUrl, fmt.Sprintf(user.UserUrl, user.SelectAllUsersQuery)))
	if err != nil {
		return object, err
	}
	return object, json.Unmarshal(body, &object)
}

// GetUsersByPage gets users by page
func (us *Uspacy) GetUsersByPage(page string) (object user.Users, err error) {
	body, err := us.doGetEmptyHeaders(us.buildURL(user.VersionUrl, fmt.Sprintf(user.UserUrl, fmt.Sprintf(user.PagePagination, page))))
	if err != nil {
		return object, err
	}
	return object, json.Unmarshal(body, &object)
}

// CreateActiveUser returns created users
func (us *Uspacy) CreateActiveUsers(entity []user.User) (object []user.User, err error) {
	body, err := us.doPostEmptyHeaders(us.buildURL(user.VersionUrl, user.CreateActiveUser), entity)
	if err != nil {
		return object, err
	}
	return object, json.Unmarshal(body, &object)
}

// PatchUser patch user by Id and return it
func (us *Uspacy) PatchUser(entity user.User) (object user.User, err error) {
	body, err := us.doPatchEmptyHeaders(us.buildURL(user.VersionUrl, fmt.Sprintf(user.UserUrl, entity.ID)), entity)
	if err != nil {
		return object, err
	}
	return object, json.Unmarshal(body, &object)
}
