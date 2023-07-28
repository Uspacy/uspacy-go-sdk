package departments

import (
	"github.com/Uspacy/uspacy-go-sdk/common"
)

const (
	VersionUrl          = "company/v1"
	DepartmentsUrl      = "departments/%v"
	DepartmentsAddUsers = "departments/%d/addUsers"
)

type (
	Departments struct {
		Data []Department `json:"data"`
		Meta common.Meta  `json:"meta"`
	}

	Department struct {
		ID                 int      `json:"id"`
		Name               string   `json:"name"`
		Description        string   `json:"description"`
		Main               bool     `json:"main"`
		HeadID             string   `json:"headId"`
		ParentDepartmentID string   `json:"parentDepartmentId"`
		UsersIds           []string `json:"usersIds"`
		Roles              []string `json:"roles"`
	}
)
