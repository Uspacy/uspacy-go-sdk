package api

import (
	"encoding/json"
	"fmt"

	"github.com/Uspacy/uspacy-go-sdk/crm"
)

// GetFields returns Fields struct for a given type of entity
func (us *Uspacy) GetFields(entity string) (fields crm.Fields, err error) {
	body, err := us.doGetEmptyHeaders(buildURL(mainHost, crm.VersionUrl, fmt.Sprintf(crm.FieldsUrl, entity, "")))
	if err != nil {
		return fields, err
	}
	return fields, json.Unmarshal(body, &fields)
}

// GetField returns Field struct for a given type of entity & field
func (us *Uspacy) GetField(entity, fieldType string) (field crm.Field, err error) {
	body, err := us.doGetEmptyHeaders(buildURL(mainHost, crm.VersionUrl, fmt.Sprintf(crm.FieldsUrl, entity, fieldType)))
	if err != nil {
		return field, err
	}
	return field, json.Unmarshal(body, &field)
}

// CreateField in CRM entity returns created field
func (us *Uspacy) CreateField(entityType crm.Entity, field interface{}) (createdField crm.Field, err error) {
	responseBody, err := us.doPostEmptyHeaders(buildURL(mainHost, crm.VersionUrl, fmt.Sprintf(crm.EntityUrl, entityType.GetUrl(crm.StatusId))), field)
	if err != nil {
		return createdField, err
	}
	return createdField, json.Unmarshal(responseBody, &createdField)
}
