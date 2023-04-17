package api

import (
	"encoding/json"
	"fmt"

	"github.com/Uspacy/uspacy-go-sdk/departments"
)

// CreateDepartment returns created department
func (us *Uspacy) CreateDepartment(departmentData departments.Department) (department departments.Department, err error) {
	body, err := us.doPostEmptyHeaders(us.buildURL(departments.VersionUrl, departments.DepartmentsUrl), departmentData)
	if err != nil {
		return department, err
	}
	return department, json.Unmarshal(body, &department)
}

// PatchDepartment patch department by Id and return it
func (us *Uspacy) PatchDepartment(departmenID int, departmentData map[string]interface{}) (department departments.Department, err error) {
	body, err := us.doPatchEmptyHeaders(us.buildURL(departments.VersionUrl, fmt.Sprintf(departments.DepartmentsUrl, departmenID)), departmentData)
	if err != nil {
		return department, err
	}
	return department, json.Unmarshal(body, &department)
}
