package api

import (
	"encoding/json"
	"fmt"

	"github.com/Uspacy/uspacy-go-sdk/crm"
)

// CreateLeads returns created contacts object
func (us *Uspacy) CreateLead(entityType crm.Entity, entity map[string]interface{}) (object crm.Leads, err error) {
	body, err := us.doPostEmptyHeaders(buildURL(mainHost, crm.VersionCRMUrl, fmt.Sprintf(crm.EntityUrl, entityType.GetUrl())), entity)
	if err != nil {
		return object, err
	}
	return object, json.Unmarshal(body, &object)
}

// CreateContact returns created contacts object
func (us *Uspacy) CreateContact(entityType crm.Entity, entity map[string]interface{}) (object crm.Contacts, err error) {
	body, err := us.doPostEmptyHeaders(buildURL(mainHost, crm.VersionCRMUrl, fmt.Sprintf(crm.EntityUrl, entityType.GetUrl())), entity)
	if err != nil {
		return object, err
	}
	return object, json.Unmarshal(body, &object)
}

// CreateCompany returns created contacts object
func (us *Uspacy) CreateCompany(entityType crm.Entity, entity map[string]interface{}) (object crm.Companies, err error) {
	body, err := us.doPostEmptyHeaders(buildURL(mainHost, crm.VersionCRMUrl, fmt.Sprintf(crm.EntityUrl, entityType.GetUrl())), entity)
	if err != nil {
		return object, err
	}
	return object, json.Unmarshal(body, &object)
}

// CreateDeals returns created contacts object
func (us *Uspacy) CreateDeal(entityType crm.Entity, entity map[string]interface{}) (object crm.Deals, err error) {
	body, err := us.doPostEmptyHeaders(buildURL(mainHost, crm.VersionCRMUrl, fmt.Sprintf(crm.EntityUrl, entityType.GetUrl())), entity)
	if err != nil {
		return object, err
	}
	return object, json.Unmarshal(body, &object)
}

// CreateFunnel returns created funnel
func (us *Uspacy) CreateFunnel(entityType crm.Entity, funnel interface{}) (object crm.Funnel, err error) {
	responseBody, err := us.doPostEmptyHeaders(buildURL(mainHost, crm.VersionCRMUrl, fmt.Sprintf(crm.FunnelUrl, entityType.GetUrl())), funnel)
	if err != nil {
		return object, err
	}
	return object, json.Unmarshal(responseBody, &object)
}

// CreateKanbanStage returns created kanban stage
func (us *Uspacy) CreateFunnelStage(entityType crm.Entity, stage interface{}) (object crm.KanbanStage, err error) {
	responseBody, err := us.doPostEmptyHeaders(buildURL(mainHost, crm.VersionCRMUrl, fmt.Sprintf(crm.KanbanStageUrl, entityType.GetUrl())), stage)
	if err != nil {
		return object, err
	}
	return object, json.Unmarshal(responseBody, &object)
}

// CreateListValues returns arrey of values for given type of CRM list
func (us *Uspacy) CreateListValues(entityType crm.Entity, listName string, lists interface{}) (object []crm.List, err error) {
	responseBody, err := us.doPostEmptyHeaders(buildURL(mainHost, crm.VersionCRMUrl, fmt.Sprintf(crm.ListsUrl, entityType.GetUrl(), listName)), lists)
	if err != nil {
		return object, err
	}
	return object, json.Unmarshal(responseBody, &object)
}

// CreateCRMField in CRM entity returns created field
func (us *Uspacy) CreateCRMField(entityType crm.Entity, body interface{}) (object crm.Field, err error) {
	responseBody, err := us.doPostEmptyHeaders(buildURL(mainHost, crm.VersionCRMUrl, fmt.Sprintf(crm.FieldsUrl, entityType.GetUrl(), "")), body)
	if err != nil {
		return object, err
	}
	return object, json.Unmarshal(responseBody, &object)
}

// GetFields returns Fields struct for a given type of entity
func (us *Uspacy) GetFields(entityType crm.Entity) (object crm.Fields, err error) {
	body, err := us.doGetEmptyHeaders(buildURL(mainHost, crm.VersionCRMUrl, fmt.Sprintf(crm.FieldsUrl, entityType.GetUrl(), "")))
	if err != nil {
		return object, err
	}
	return object, json.Unmarshal(body, &object)
}

// GetField returns Field struct for a given type of entity & field
func (us *Uspacy) GetField(entityType crm.Entity, fieldType string) (object crm.Field, err error) {
	body, err := us.doGetEmptyHeaders(buildURL(mainHost, crm.VersionCRMUrl, fmt.Sprintf(crm.FieldsUrl, entityType.GetUrl(), fieldType)))
	if err != nil {
		return object, err
	}
	return object, json.Unmarshal(body, &object)
}
