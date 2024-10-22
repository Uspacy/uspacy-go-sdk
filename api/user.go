package api

import (
	"encoding/json"
	"fmt"

	"github.com/Uspacy/uspacy-go-sdk/user"
)

// GetAllUsers gets all users
func (us *Uspacy) GetAllUsers() (users []user.User, err error) {
	body, err := us.doGetEmptyHeaders(us.buildURL(user.VersionUrl, fmt.Sprintf(user.UserUrl, user.SelectAllUsersQuery)))
	if err != nil {
		return users, err
	}
	return users, json.Unmarshal(body, &users)
}

// GetUsersByPage gets users by page
func (us *Uspacy) GetUsersByPage(page string) (users user.Users, err error) {
	body, err := us.doGetEmptyHeaders(us.buildURL(user.VersionUrl, fmt.Sprintf(user.UserUrl, fmt.Sprintf(user.PagePagination, page))))
	if err != nil {
		return users, err
	}
	return users, json.Unmarshal(body, &users)
}

// CreateActiveUser returns created users
func (us *Uspacy) CreateActiveUsers(usersData []user.UsersInvite, headers ...map[string]string) (users []user.CreatedActiveUser, err error) {
	body, _, err := us.doPost(us.buildURL(user.VersionUrl, user.CreateActiveUser), usersData, headers...)
	if err != nil {
		return users, err
	}
	return users, json.Unmarshal(body, &users)
}

// PatchUser patch user by Id and return it
func (us *Uspacy) PatchUser(userData user.User) (_user user.User, err error) {
	body, err := us.doPatchEmptyHeaders(us.buildURL(user.VersionUrl, fmt.Sprintf(user.UserUrl, userData.ID)), userData)
	if err != nil {
		return _user, err
	}
	return _user, json.Unmarshal(body, &_user)
}
