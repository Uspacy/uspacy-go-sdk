package api

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Uspacy/uspacy-go-sdk/v2/departments"
)

// GetDepartments returns list of departments
func (us *Uspacy) GetDepartments(ctx context.Context) (departmentsArrey departments.Departments, err error) {
	body, err := us.doGetEmptyHeaders(ctx, us.buildURL(departments.VersionUrl, fmt.Sprintf(departments.DepartmentsUrl, "")))
	if err != nil {
		return departmentsArrey, err
	}
	return departmentsArrey, json.Unmarshal(body, &departmentsArrey)
}

// CreateDepartment returns created department
func (us *Uspacy) CreateDepartment(ctx context.Context, departmentData departments.Department, opts ...RequestOption) (department departments.Department, err error) {
	body, _, err := us.doPost(ctx, us.buildURL(departments.VersionUrl, fmt.Sprintf(departments.DepartmentsUrl, "")), departmentData, opts...)
	if err != nil {
		return department, err
	}
	return department, json.Unmarshal(body, &department)
}

// PatchDepartment patch department by Id and return it
func (us *Uspacy) PatchDepartment(ctx context.Context, departmentID int, departmentData map[string]any) (department departments.Department, err error) {
	body, err := us.doPatchEmptyHeaders(ctx, us.buildURL(departments.VersionUrl, fmt.Sprintf(departments.DepartmentsUrl, departmentID)), departmentData)
	if err != nil {
		return department, err
	}
	return department, json.Unmarshal(body, &department)
}

// DepartmentAddUsers patch department by Id and return it
func (us *Uspacy) DepartmentAddUsers(ctx context.Context, departmentID int, usersIds []int) (department departments.Department, err error) {
	body, err := us.doPatchEmptyHeaders(ctx, us.buildURL(departments.VersionUrl, fmt.Sprintf(departments.DepartmentsAddUsers, departmentID)), usersIds)
	if err != nil {
		return department, err
	}
	return department, json.Unmarshal(body, &department)
}
