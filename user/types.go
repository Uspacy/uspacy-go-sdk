package user

import (
	"github.com/Uspacy/uspacy-go-sdk/common"
)

const (
	VersionUrl       = "company/v1/"
	CreateActiveUser = "invites/email/import_registered"
	UserUrl          = "users/%s"
	PagePagination   = "?page=%s"
	AllUsersParametr = "?show=all&list=all"
)

type Users struct {
	Data []User      `json:"data"`
	Meta common.Meta `json:"meta"`
}

type User struct {
	ID               int                  `json:"id"`
	Active           bool                 `json:"active"`
	FirstName        string               `json:"firstName"`
	LastName         string               `json:"lastName"`
	Patronymic       string               `json:"patronymic"`
	Birthday         int                  `json:"birthday"`
	ShowBirthYear    bool                 `json:"showBirthYear"`
	AboutMyself      string               `json:"aboutMyself"`
	Position         string               `json:"position"`
	Specialization   string               `json:"specialization"`
	Country          string               `json:"country"`
	City             string               `json:"city"`
	Avatar           string               `json:"avatar"`
	EmailInvitation  string               `json:"emailInvitation"`
	DateOfInvitation string               `json:"dateOfInvitation"`
	Registered       bool                 `json:"registered"`
	Email            []common.ContactData `json:"email"`
	Phone            []common.ContactData `json:"phone"`
	SocialMedia      []common.SocialMedia `json:"socialMedia"`
	DepartmentsIds   []string             `json:"departmentsIds"`
	Roles            []string             `json:"roles"`
}
