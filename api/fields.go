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

// CreateField returns created field
func (us *Uspacy) CreateField(statusID crm.Entity, field interface{}) (crm.Field, error) {
	var createdField crm.Field
	responseBody, err := us.doPostEmptyHeaders(buildURL(mainHost, crm.VersionUrl, fmt.Sprintf(crm.FunnelUrl, statusID.GetUrl(crm.StatusId))), field)
	if err != nil {
		return createdField, err
	}
	return createdField, json.Unmarshal(responseBody, &createdField)
}

// CreateField returns arrey of created lists
func (us *Uspacy) CreateList(entityType, listValue crm.Entity, field interface{}) ([]crm.List, error) {
	var createdLists []crm.List
	responseBody, err := us.doPostEmptyHeaders(
		buildURL(
			mainHost,
			crm.VersionUrl,
			fmt.Sprintf(crm.ListsUrl, entityType.GetUrl(crm.EntityType), listValue.GetUrl(crm.ListValue))),
		field)
	if err != nil {
		return createdLists, err
	}
	return createdLists, json.Unmarshal(responseBody, &createdLists)
}
