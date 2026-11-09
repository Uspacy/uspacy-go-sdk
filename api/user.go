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

// GetUser returns the user with the given ID.
func (us *Uspacy) GetUser(ctx context.Context, userID int, opts ...RequestOption) (user.User, error) {
	body, err := us.doGet(ctx, us.buildURL(user.VersionUrl, fmt.Sprintf(user.UserUrl, userID)), opts...)
	return decodeJSON[user.User](body, err)
}

// CreateActiveUsers imports users as already registered (invites/email/import_registered)
// and returns the created users.
func (us *Uspacy) CreateActiveUsers(ctx context.Context, usersData []user.UsersInvite, opts ...RequestOption) ([]user.CreatedActiveUser, error) {
	body, _, err := us.doPost(ctx, us.buildURL(user.VersionUrl, user.CreateActiveUser), usersData, opts...)
	return decodeJSON[[]user.CreatedActiveUser](body, err)
}

// PatchUser updates a user and returns it. Only the keys of userData are sent; the other
// fields keep their values.
func (us *Uspacy) PatchUser(ctx context.Context, userID int, userData map[string]any, opts ...RequestOption) (user.User, error) {
	body, err := us.doPatch(ctx, us.buildURL(user.VersionUrl, fmt.Sprintf(user.UserUrl, userID)), userData, opts...)
	return decodeJSON[user.User](body, err)
}

// GetUserFields returns the fields of users, the system ones included. A list or label
// field comes with its values.
func (us *Uspacy) GetUserFields(ctx context.Context, opts ...RequestOption) ([]user.Field, error) {
	body, err := us.doGet(ctx, us.buildURL(user.VersionUrl, fmt.Sprintf(user.FieldsUrl, "")), opts...)
	return decodeJSON[[]user.Field](body, err)
}

// GetUserField returns the field of users with the given code.
func (us *Uspacy) GetUserField(ctx context.Context, code string, opts ...RequestOption) (user.Field, error) {
	body, err := us.doGet(ctx, us.buildURL(user.VersionUrl, fmt.Sprintf(user.FieldsUrl, code)), opts...)
	return decodeJSON[user.Field](body, err)
}

// CreateUserField creates a custom field of users and returns it. fieldData is a
// user.Field or a map; name and type are required, and the code is generated from the
// name when empty. Keys left out of a map get the backend defaults (editable and show
// true). The values of a list field are not saved here: add them with
// CreateUserListValues.
func (us *Uspacy) CreateUserField(ctx context.Context, fieldData any, opts ...RequestOption) (user.Field, error) {
	body, _, err := us.doPost(ctx, us.buildURL(user.VersionUrl, fmt.Sprintf(user.FieldsUrl, "")), fieldData, opts...)
	return decodeJSON[user.Field](body, err)
}

// PatchUserField updates the field of users with the given code and returns it. Only the
// keys of fieldData are sent; the other properties keep their values.
func (us *Uspacy) PatchUserField(ctx context.Context, code string, fieldData map[string]any, opts ...RequestOption) (user.Field, error) {
	body, err := us.doPatch(ctx, us.buildURL(user.VersionUrl, fmt.Sprintf(user.FieldsUrl, code)), fieldData, opts...)
	return decodeJSON[user.Field](body, err)
}

// DeleteUserField deletes the custom field of users with the given code, with its values.
func (us *Uspacy) DeleteUserField(ctx context.Context, code string, opts ...RequestOption) error {
	_, err := us.doDelete(ctx, us.buildURL(user.VersionUrl, fmt.Sprintf(user.FieldsUrl, code)), nil, opts...)
	return err
}

// GetUserListValues returns the values of the list field of users with the given code.
func (us *Uspacy) GetUserListValues(ctx context.Context, code string, opts ...RequestOption) ([]user.ListValue, error) {
	body, err := us.doGet(ctx, us.buildURL(user.VersionUrl, fmt.Sprintf(user.ListsUrl, code)), opts...)
	return decodeJSON[[]user.ListValue](body, err)
}

// CreateUserListValues adds values to the list field of users with the given code and
// returns all of the field's values. A value that already exists (same Value) is updated.
func (us *Uspacy) CreateUserListValues(ctx context.Context, code string, values []user.ListValue, opts ...RequestOption) ([]user.ListValue, error) {
	body, _, err := us.doPost(ctx, us.buildURL(user.VersionUrl, fmt.Sprintf(user.ListsUrl, code)), values, opts...)
	return decodeJSON[[]user.ListValue](body, err)
}

// DeleteUserListValue deletes the value from the list field of users with the given code.
func (us *Uspacy) DeleteUserListValue(ctx context.Context, code, value string, opts ...RequestOption) error {
	_, err := us.doDelete(ctx, us.buildURL(user.VersionUrl, fmt.Sprintf(user.ListValueUrl, code, value)), nil, opts...)
	return err
}
