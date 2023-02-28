package api

import (
	"encoding/json"
	"fmt"
	"uspacy-go-sdk/crm"
)

// GetFields returns Fields struct for a given type of field
func (us *Uspacy) GetFields(field string) (crm.Fields, error) {
	var fields crm.Fields

	body, err := us.doGetEmptyHeaders(buildURL(mainHost, crm.VersionUrl, fmt.Sprintf(crm.GetFieldsUrl, field)))
	if err != nil {
		return fields, err
	}
	return fields, json.Unmarshal(body, &fields)
}
