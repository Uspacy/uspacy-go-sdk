package api

import (
	"context"
	"fmt"

	"github.com/Uspacy/uspacy-go-sdk/v2/user"
)

// GetAllUsers returns all users.
func (us *Uspacy) GetAllUsers(ctx context.Context, opts ...RequestOption) ([]user.User, error) {
	body, err := us.doGet(ctx, us.buildURL(user.VersionUrl, fmt.Sprintf(user.UserUrl, user.SelectAllUsersQuery)), opts...)
	return decodeJSON[[]user.User](body, err)
}

// GetUsersByPage returns one page of users.
func (us *Uspacy) GetUsersByPage(ctx context.Context, page string, opts ...RequestOption) (user.Users, error) {
	body, err := us.doGet(ctx, us.buildURL(user.VersionUrl, fmt.Sprintf(user.UserUrl, fmt.Sprintf(user.PagePagination, page))), opts...)
	return decodeJSON[user.Users](body, err)
}

// CreateActiveUsers imports users as already registered (invites/email/import_registered)
// and returns the created users.
func (us *Uspacy) CreateActiveUsers(ctx context.Context, usersData []user.UsersInvite, opts ...RequestOption) ([]user.CreatedActiveUser, error) {
	body, _, err := us.doPost(ctx, us.buildURL(user.VersionUrl, user.CreateActiveUser), usersData, opts...)
	return decodeJSON[[]user.CreatedActiveUser](body, err)
}

// PatchUser updates a user and returns it.
func (us *Uspacy) PatchUser(ctx context.Context, userData user.User, opts ...RequestOption) (user.User, error) {
	body, err := us.doPatch(ctx, us.buildURL(user.VersionUrl, fmt.Sprintf(user.UserUrl, userData.ID)), userData, opts...)
	return decodeJSON[user.User](body, err)
}
