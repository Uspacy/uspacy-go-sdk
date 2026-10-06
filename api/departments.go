package api

import (
	"context"
	"fmt"

	"github.com/Uspacy/uspacy-go-sdk/v2/departments"
)

// GetDepartments returns list of departments
func (us *Uspacy) GetDepartments(ctx context.Context) (departments.Departments, error) {
	body, err := us.doGetEmptyHeaders(ctx, us.buildURL(departments.VersionUrl, fmt.Sprintf(departments.DepartmentsUrl, "")))
	return decodeJSON[departments.Departments](body, err)
}

// CreateDepartment returns created department
func (us *Uspacy) CreateDepartment(ctx context.Context, departmentData departments.Department, opts ...RequestOption) (departments.Department, error) {
	body, _, err := us.doPost(ctx, us.buildURL(departments.VersionUrl, fmt.Sprintf(departments.DepartmentsUrl, "")), departmentData, opts...)
	return decodeJSON[departments.Department](body, err)
}

// PatchDepartment patch department by Id and return it
func (us *Uspacy) PatchDepartment(ctx context.Context, departmentID int, departmentData map[string]any) (departments.Department, error) {
	body, err := us.doPatchEmptyHeaders(ctx, us.buildURL(departments.VersionUrl, fmt.Sprintf(departments.DepartmentsUrl, departmentID)), departmentData)
	return decodeJSON[departments.Department](body, err)
}

// DepartmentAddUsers patch department by Id and return it
func (us *Uspacy) DepartmentAddUsers(ctx context.Context, departmentID int, usersIds []int) (departments.Department, error) {
	body, err := us.doPatchEmptyHeaders(ctx, us.buildURL(departments.VersionUrl, fmt.Sprintf(departments.DepartmentsAddUsers, departmentID)), usersIds)
	return decodeJSON[departments.Department](body, err)
}
