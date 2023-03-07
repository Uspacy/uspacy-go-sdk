package api

import (
	"encoding/json"
	"fmt"
	"uspacy-go-sdk/crm"
)

// GetFields returns Fields struct for a given type of entity
func (us *Uspacy) GetFields(entity string) (crm.Fields, error) {
	var fields crm.Fields

	body, err := us.doGetEmptyHeaders(buildURL(mainHost, crm.VersionUrl, fmt.Sprintf(crm.GetFieldsUrl, entity, "")))
	if err != nil {
		return fields, err
	}
	return fields, json.Unmarshal(body, &fields)
}

// GetField returns Field struct for a given type of entity & field
func (us *Uspacy) GetField(entity, field_type string) (crm.Field, error) {
	var field crm.Field

	body, err := us.doGetEmptyHeaders(buildURL(mainHost, crm.VersionUrl, fmt.Sprintf(crm.GetFieldsUrl, entity, field_type)))
	if err != nil {
		return field, err
	}
	return field, json.Unmarshal(body, &field)
}
