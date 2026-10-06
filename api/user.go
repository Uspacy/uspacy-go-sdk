package api

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Uspacy/uspacy-go-sdk/v2/user"
)

// GetAllUsers gets all users
func (us *Uspacy) GetAllUsers(ctx context.Context) (users []user.User, err error) {
	body, err := us.doGetEmptyHeaders(ctx, us.buildURL(user.VersionUrl, fmt.Sprintf(user.UserUrl, user.SelectAllUsersQuery)))
	if err != nil {
		return users, err
	}
	return users, json.Unmarshal(body, &users)
}

// GetUsersByPage gets users by page
func (us *Uspacy) GetUsersByPage(ctx context.Context, page string) (users user.Users, err error) {
	body, err := us.doGetEmptyHeaders(ctx, us.buildURL(user.VersionUrl, fmt.Sprintf(user.UserUrl, fmt.Sprintf(user.PagePagination, page))))
	if err != nil {
		return users, err
	}
	return users, json.Unmarshal(body, &users)
}

// CreateActiveUsers returns created users
func (us *Uspacy) CreateActiveUsers(ctx context.Context, usersData []user.UsersInvite, opts ...RequestOption) (users []user.CreatedActiveUser, err error) {
	body, _, err := us.doPost(ctx, us.buildURL(user.VersionUrl, user.CreateActiveUser), usersData, opts...)
	if err != nil {
		return users, err
	}
	return users, json.Unmarshal(body, &users)
}

// PatchUser patch user by Id and return it
func (us *Uspacy) PatchUser(ctx context.Context, userData user.User) (_user user.User, err error) {
	body, err := us.doPatchEmptyHeaders(ctx, us.buildURL(user.VersionUrl, fmt.Sprintf(user.UserUrl, userData.ID)), userData)
	if err != nil {
		return _user, err
	}
	return _user, json.Unmarshal(body, &_user)
}
