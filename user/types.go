package user

import (
	"encoding/json"

	"github.com/Uspacy/uspacy-go-sdk/v2/common"
)

const (
	VersionUrl          = "company/v1"
	CreateActiveUser    = "invites/email/import_registered"
	UserUrl             = "users/%v"
	PagePagination      = "?page=%s"
	SelectAllUsersQuery = "?show=all&list=all"
	FieldsUrl           = "custom_fields/users/fields/%s"
	ListsUrl            = "custom_fields/users/lists/%s"
	ListValueUrl        = "custom_fields/users/lists/%s/%s"
)

type Users struct {
	Data []User      `json:"data"`
	Meta common.Meta `json:"meta"`
}

type UsersInvite struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Password  string `json:"password"`
	Domain    string `json:"domain"`
}

type CreatedActiveUser struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

type User struct {
	ID               int                  `json:"id"`
	Active           bool                 `json:"active"`
	FirstName        string               `json:"firstName"`
	LastName         string               `json:"lastName"`
	Patronymic       string               `json:"patronymic"`
	Birthday         int                  `json:"birthday"`
	ShowBirthYear    bool                 `json:"showBirthYear"`
	AboutMyself      string               `json:"about_myself"`
	Position         string               `json:"position"`
	Specialization   string               `json:"specialization"`
	Country          string               `json:"country"`
	City             string               `json:"city"`
	Avatar           string               `json:"avatar"`
	EmailInvitation  any                  `json:"emailInvitation"`
	DateOfInvitation any                  `json:"dateOfInvitation"`
	Registered       bool                 `json:"registered"`
	Email            []common.ContactData `json:"email"`
	Phone            []common.ContactData `json:"phone"`
	SocialMedia      []common.SocialMedia `json:"socialMedia"`
	DepartmentsIds   any                  `json:"departmentsIds"`
	Roles            any                  `json:"roles"`

	// Raw holds every key of the user as received, custom fields included under their
	// code. Read a custom field with FieldValue.
	Raw map[string]json.RawMessage `json:"-"`
}

// UnmarshalJSON decodes the user into its fields and keeps every key in Raw.
func (u *User) UnmarshalJSON(data []byte) error {
	type plain User // without this method, so decoding it does not recurse
	if err := json.Unmarshal(data, (*plain)(u)); err != nil {
		return err
	}
	// A new map, not &u.Raw: decoding into a non-nil map adds to it, so a reused User would
	// keep the keys of an earlier decode and FieldValue would report values it does not have.
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	u.Raw = raw
	return nil
}

// FieldValue decodes the value of the field with the given code (see GetUserFields) into
// v and reports whether the user has one. It is false when the key is absent (no such
// field, or one the current user may not view) and when the value is null or "", which
// the backend sends for an empty list field.
func (u User) FieldValue(code string, v any) (bool, error) {
	raw, ok := u.Raw[code]
	if !ok || string(raw) == "null" || string(raw) == `""` {
		return false, nil
	}
	if err := json.Unmarshal(raw, v); err != nil {
		return false, err
	}
	return true, nil
}

// Field is a field of users: a system one (first_name, email, ...) or a custom one.
// The backend sends "" for an unset FieldSectionId, EntityReferenceId and Sort; they decode
// to 0 and, when 0, are left out of a request body.
type Field struct {
	Name              string            `json:"name"`
	Code              string            `json:"code"`
	Type              string            `json:"type"`
	Required          bool              `json:"required"`
	Editable          bool              `json:"editable"`
	Show              bool              `json:"show"`
	Hidden            bool              `json:"hidden"`
	Multiple          bool              `json:"multiple"`
	SystemField       bool              `json:"system_field"`
	BaseField         bool              `json:"base_field"`
	FieldSectionId    common.FlexInt    `json:"field_section_id,omitempty"`
	EntityReferenceId common.FlexInt    `json:"entity_reference_id,omitempty"`
	Sort              common.FlexInt    `json:"sort,omitempty"`
	DefaultValue      string            `json:"default_value"`
	Tooltip           string            `json:"tooltip"`
	DisplayInTabs     bool              `json:"display_in_tabs"`
	Permissions       *FieldPermissions `json:"permissions"`      // nil = everyone may view and edit; only admins see it
	Values            []ListValue       `json:"values,omitempty"` // list and label fields only
	Access            *FieldAccess      `json:"access,omitempty"` // what the current user may do; ignored in requests
}

// FieldPermissions limits who may view and edit the values of a custom field.
type FieldPermissions struct {
	View FieldPermission `json:"view"`
	Edit FieldPermission `json:"edit"`
}

// FieldPermission grants an action to everyone (All) or only to the given roles and users.
type FieldPermission struct {
	All   bool  `json:"all"`
	Roles []int `json:"roles,omitempty"`
	Users []int `json:"users,omitempty"`
}

// FieldAccess says whether the current user may view and edit the values of a field.
type FieldAccess struct {
	View bool `json:"view"`
	Edit bool `json:"edit"`
}

// ListValue is a value of a list or label field.
type ListValue struct {
	Title    string         `json:"title"`
	Value    string         `json:"value"`
	Color    string         `json:"color"`
	Sort     common.FlexInt `json:"sort"`
	Selected bool           `json:"selected"`
}
