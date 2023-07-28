package api

import (
	"encoding/json"
	"fmt"

	"github.com/Uspacy/uspacy-go-sdk/departments"
)

// GetDepartments returns list of departments
func (us *Uspacy) GetDepartments() (departmentsArrey departments.Departments, err error) {
	body, err := us.doGetEmptyHeaders(us.buildURL(departments.VersionUrl, fmt.Sprintf(departments.DepartmentsUrl, "")))
	if err != nil {
		return departmentsArrey, err
	}
	return departmentsArrey, json.Unmarshal(body, &departmentsArrey)
}

// CreateDepartment returns created department
func (us *Uspacy) CreateDepartment(departmentData departments.Department) (department departments.Department, err error) {
	body, err := us.doPostEmptyHeaders(us.buildURL(departments.VersionUrl, fmt.Sprintf(departments.DepartmentsUrl, "")), departmentData)
	if err != nil {
		return department, err
	}
	return department, json.Unmarshal(body, &department)
}

// PatchDepartment patch department by Id and return it
func (us *Uspacy) PatchDepartment(departmentID int, departmentData map[string]interface{}) (department departments.Department, err error) {
	body, err := us.doPatchEmptyHeaders(us.buildURL(departments.VersionUrl, fmt.Sprintf(departments.DepartmentsUrl, departmentID)), departmentData)
	if err != nil {
		return department, err
	}
	return department, json.Unmarshal(body, &department)
}

// DepartmentAddUsers patch department by Id and return it
func (us *Uspacy) DepartmentAddUsers(departmentID int, usersIds []int) (department departments.Department, err error) {
	body, err := us.doPatchEmptyHeaders(us.buildURL(departments.VersionUrl, fmt.Sprintf(departments.DepartmentsAddUsers, departmentID)), usersIds)
	if err != nil {
		return department, err
	}
	return department, json.Unmarshal(body, &department)
}
