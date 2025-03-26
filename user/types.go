package user

import (
	"github.com/Uspacy/uspacy-go-sdk/common"
)

const (
	VersionUrl          = "company/v1"
	CreateActiveUser    = "invites/email/import_registered"
	UserUrl             = "users/%v"
	PagePagination      = "?page=%s"
	SelectAllUsersQuery = "?show=all&list=all"
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
	AboutMyself      string               `json:"aboutMyself"`
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
}
