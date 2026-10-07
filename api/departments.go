package api

import (
	"context"
	"fmt"

	"github.com/Uspacy/uspacy-go-sdk/v2/departments"
)

// GetDepartments returns all departments.
func (us *Uspacy) GetDepartments(ctx context.Context, opts ...RequestOption) (departments.Departments, error) {
	body, err := us.doGet(ctx, us.buildURL(departments.VersionUrl, fmt.Sprintf(departments.DepartmentsUrl, "")), opts...)
	return decodeJSON[departments.Departments](body, err)
}

// CreateDepartment creates a department and returns it.
func (us *Uspacy) CreateDepartment(ctx context.Context, departmentData departments.Department, opts ...RequestOption) (departments.Department, error) {
	body, _, err := us.doPost(ctx, us.buildURL(departments.VersionUrl, fmt.Sprintf(departments.DepartmentsUrl, "")), departmentData, opts...)
	return decodeJSON[departments.Department](body, err)
}

// PatchDepartment updates a department and returns it.
func (us *Uspacy) PatchDepartment(ctx context.Context, departmentID int, departmentData map[string]any, opts ...RequestOption) (departments.Department, error) {
	body, err := us.doPatch(ctx, us.buildURL(departments.VersionUrl, fmt.Sprintf(departments.DepartmentsUrl, departmentID)), departmentData, opts...)
	return decodeJSON[departments.Department](body, err)
}

// DepartmentAddUsers adds users to a department and returns the department.
func (us *Uspacy) DepartmentAddUsers(ctx context.Context, departmentID int, usersIds []int, opts ...RequestOption) (departments.Department, error) {
	body, err := us.doPatch(ctx, us.buildURL(departments.VersionUrl, fmt.Sprintf(departments.DepartmentsAddUsers, departmentID)), usersIds, opts...)
	return decodeJSON[departments.Department](body, err)
}
