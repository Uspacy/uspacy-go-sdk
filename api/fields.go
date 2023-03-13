package api

import (
	"encoding/json"
	"fmt"

	"github.com/Uspacy/uspacy-go-sdk/crm"
)

// GetFields returns Fields struct for a given type of entity
func (us *Uspacy) GetFields(entityType crm.Entity) (fields crm.Fields, err error) {
	body, err := us.doGetEmptyHeaders(buildURL(mainHost, crm.VersionCRMUrl, fmt.Sprintf(crm.FieldsUrl, entityType.GetUrl(crm.EntityType), "")))
	if err != nil {
		return fields, err
	}
	return fields, json.Unmarshal(body, &fields)
}

// GetField returns Field struct for a given type of entity & field
func (us *Uspacy) GetField(entityType crm.Entity, fieldType string) (field crm.Field, err error) {
	body, err := us.doGetEmptyHeaders(buildURL(mainHost, crm.VersionCRMUrl, fmt.Sprintf(crm.FieldsUrl, entityType.GetUrl(crm.EntityType), fieldType)))
	if err != nil {
		return field, err
	}
	return field, json.Unmarshal(body, &field)
}

// CreateCRMField in CRM entity returns created field
func (us *Uspacy) CreateCRMField(entityType crm.Entity, field interface{}) (createdField crm.Field, err error) {
	responseBody, err := us.doPostEmptyHeaders(buildURL(mainHost, crm.VersionCRMUrl, fmt.Sprintf(crm.FieldsUrl, entityType.GetUrl(crm.EntityType), "")), field)
	if err != nil {
		return createdField, err
	}
	return createdField, json.Unmarshal(responseBody, &createdField)
}
