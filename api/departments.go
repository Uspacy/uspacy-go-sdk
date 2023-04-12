package api

import (
	"encoding/json"
	"fmt"

	"github.com/Uspacy/uspacy-go-sdk/departments"
)

// CreateActiveDepartment returns created department
func (us *Uspacy) CreateActiveDepartment(departmentData departments.Department) (department []departments.Department, err error) {
	body, err := us.doPostEmptyHeaders(us.buildURL(departments.VersionUrl, departments.DepartmentsUrl), departmentData)
	if err != nil {
		return department, err
	}
	return department, json.Unmarshal(body, &department)
}

// PatchDepartment patch user by Id and return it
func (us *Uspacy) PatchDepartment(departmentData departments.Department) (department []departments.Department, err error) {
	body, err := us.doPatchEmptyHeaders(us.buildURL(departments.VersionUrl, fmt.Sprintf(departments.DepartmentsUrl, departmentData.ID)), departmentData)
	if err != nil {
		return department, err
	}
	return department, json.Unmarshal(body, &department)
}
