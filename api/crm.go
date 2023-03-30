package api

import (
	"encoding/json"
	"fmt"

	"github.com/Uspacy/uspacy-go-sdk/crm"
)

// CreateObject this method does not return any object, just error
func (us *Uspacy) CreateObject(service Service, entityType crm.Entity, entity map[string]interface{}) error {
	_, err := us.doPostEmptyHeaders(us.buildURL(service.getService(), fmt.Sprintf(crm.EntityUrl, entityType.GetUrl())), entity)
	if err != nil {
		return err
	}
	return nil

}

// CreateContact returns created contact object
func (us *Uspacy) CreateContact(entityType crm.Entity, entity map[string]interface{}) (object crm.Contacts, err error) {
	body, err := us.doPostEmptyHeaders(us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.EntityUrl, entityType.GetUrl())), entity)
	if err != nil {
		return object, err
	}
	return object, json.Unmarshal(body, &object)
}

// CreateCompany returns created company object
func (us *Uspacy) CreateCompany(entityType crm.Entity, entity map[string]interface{}) (object crm.Companies, err error) {
	body, err := us.doPostEmptyHeaders(us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.EntityUrl, entityType.GetUrl())), entity)
	if err != nil {
		return object, err
	}
	return object, json.Unmarshal(body, &object)
}

// CreateLeads returns created lead object
func (us *Uspacy) CreateLead(entityType crm.Entity, entity map[string]interface{}) (object crm.Leads, err error) {
	body, err := us.doPostEmptyHeaders(us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.EntityUrl, entityType.GetUrl())), entity)
	if err != nil {
		return object, err
	}
	return object, json.Unmarshal(body, &object)
}

// CreateDeals returns created deal object
func (us *Uspacy) CreateDeal(entityType crm.Entity, entity map[string]interface{}) (object crm.Deals, err error) {
	body, err := us.doPostEmptyHeaders(us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.EntityUrl, entityType.GetUrl())), entity)
	if err != nil {
		return object, err
	}
	return object, json.Unmarshal(body, &object)
}

// GetField returns Field struct for a given type of entity & field
func (us *Uspacy) GetField(entityType crm.Entity, fieldType string) (object crm.Field, err error) {
	body, err := us.doGetEmptyHeaders(us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.FieldsUrl, entityType.GetUrl(), fieldType)))
	if err != nil {
		return object, err
	}
	return object, json.Unmarshal(body, &object)
}

// GetFields returns Fields struct for a given type of entity
func (us *Uspacy) GetFields(entityType crm.Entity) (object crm.Fields, err error) {
	body, err := us.doGetEmptyHeaders(us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.FieldsUrl, entityType.GetUrl(), "")))
	if err != nil {
		return object, err
	}
	return object, json.Unmarshal(body, &object)
}

// CreateFunnel returns created funnel
func (us *Uspacy) CreateFunnel(entityType crm.Entity, funnel interface{}) (object crm.Funnel, err error) {
	responseBody, err := us.doPostEmptyHeaders(us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.FunnelUrl, entityType.GetUrl())), funnel)
	if err != nil {
		return object, err
	}
	return object, json.Unmarshal(responseBody, &object)
}

// CreateFunnelStage returns created kanban stage
func (us *Uspacy) CreateFunnelStage(entityType crm.Entity, stage interface{}) (object crm.KanbanStage, err error) {
	responseBody, err := us.doPostEmptyHeaders(us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.KanbanStageUrl, entityType.GetUrl())), stage)
	if err != nil {
		return object, err
	}
	return object, json.Unmarshal(responseBody, &object)
}

// Move a funnel stage
func (us *Uspacy) MoveFunnelStage(entityType crm.Entity, entityId int64, stageId string) (err error) {
	_, err = us.doPostEmptyHeaders(us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.MoveKanbanStageUrl, entityType.GetUrl(), entityId, stageId)), nil)
	return err
}

// CreateCRMField in CRM entity returns created field
func (us *Uspacy) CreateCRMField(entityType crm.Entity, body interface{}) (object crm.Field, err error) {
	responseBody, err := us.doPostEmptyHeaders(us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.FieldsUrl, entityType.GetUrl(), "")), body)
	if err != nil {
		return object, err
	}
	return object, json.Unmarshal(responseBody, &object)
}

// CreateListValues returns arrey of values for given type of CRM list
func (us *Uspacy) CreateListValues(entityType crm.Entity, listName string, lists interface{}) (object []crm.List, err error) {
	responseBody, err := us.doPostEmptyHeaders(us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.ListsUrl, entityType.GetUrl(), listName)), lists)
	if err != nil {
		return object, err
	}
	return object, json.Unmarshal(responseBody, &object)
}
