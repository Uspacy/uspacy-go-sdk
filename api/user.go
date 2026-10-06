package api

import (
	"context"
	"fmt"

	"github.com/Uspacy/uspacy-go-sdk/v2/user"
)

// GetAllUsers gets all users
func (us *Uspacy) GetAllUsers(ctx context.Context) ([]user.User, error) {
	body, err := us.doGetEmptyHeaders(ctx, us.buildURL(user.VersionUrl, fmt.Sprintf(user.UserUrl, user.SelectAllUsersQuery)))
	return decodeJSON[[]user.User](body, err)
}

// GetUsersByPage gets users by page
func (us *Uspacy) GetUsersByPage(ctx context.Context, page string) (user.Users, error) {
	body, err := us.doGetEmptyHeaders(ctx, us.buildURL(user.VersionUrl, fmt.Sprintf(user.UserUrl, fmt.Sprintf(user.PagePagination, page))))
	return decodeJSON[user.Users](body, err)
}

// CreateActiveUsers returns created users
func (us *Uspacy) CreateActiveUsers(ctx context.Context, usersData []user.UsersInvite, opts ...RequestOption) ([]user.CreatedActiveUser, error) {
	body, _, err := us.doPost(ctx, us.buildURL(user.VersionUrl, user.CreateActiveUser), usersData, opts...)
	return decodeJSON[[]user.CreatedActiveUser](body, err)
}

// PatchUser patch user by Id and return it
func (us *Uspacy) PatchUser(ctx context.Context, userData user.User) (user.User, error) {
	body, err := us.doPatchEmptyHeaders(ctx, us.buildURL(user.VersionUrl, fmt.Sprintf(user.UserUrl, userData.ID)), userData)
	return decodeJSON[user.User](body, err)
}
