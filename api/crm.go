package api

import (
	"encoding/json"
	"fmt"

	"github.com/Uspacy/uspacy-go-sdk/crm"
)

// CreateEntity this method does not return any object, just error
func (us *Uspacy) CreateEntity(entityType crm.Entity, entityData map[string]interface{}) error {
	_, err := us.doPostEmptyHeaders(us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.EntityUrl, entityType.GetUrl())), entityData)
	if err != nil {
		return err
	}
	return nil
}

// PatchEntity this method does not return any object, just error
func (us *Uspacy) PatchEntity(entityType crm.Entity, id string, entityData map[string]interface{}) error {
	_, err := us.doPatchEmptyHeaders(us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.EntityUrl, entityType.GetUrl()))+id, entityData)
	if err != nil {
		return err
	}
	return nil
}

// CreateContact returns created contact object
func (us *Uspacy) CreateContact(entityType crm.Entity, contactData map[string]interface{}) (contact crm.Contact, err error) {
	body, err := us.doPostEmptyHeaders(us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.EntityUrl, entityType.GetUrl())), contactData)
	if err != nil {
		return contact, err
	}
	return contact, json.Unmarshal(body, &contact)
}

// CreateCompany returns created company object
func (us *Uspacy) CreateCompany(entityType crm.Entity, companyData map[string]interface{}) (company crm.Company, err error) {
	body, err := us.doPostEmptyHeaders(us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.EntityUrl, entityType.GetUrl())), companyData)
	if err != nil {
		return company, err
	}
	return company, json.Unmarshal(body, &company)
}

// CreateLeads returns created lead object
func (us *Uspacy) CreateLead(entityType crm.Entity, leadData map[string]interface{}) (lead crm.Lead, err error) {
	body, err := us.doPostEmptyHeaders(us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.EntityUrl, entityType.GetUrl())), leadData)
	if err != nil {
		return lead, err
	}
	return lead, json.Unmarshal(body, &lead)
}

// CreateDeals returns created deal object
func (us *Uspacy) CreateDeal(entityType crm.Entity, dealData map[string]interface{}) (deal crm.Deal, err error) {
	body, err := us.doPostEmptyHeaders(us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.EntityUrl, entityType.GetUrl())), dealData)
	if err != nil {
		return deal, err
	}
	return deal, json.Unmarshal(body, &deal)
}

// CreateTask returns created task object
func (us *Uspacy) CreateTaskCRM(entityType crm.Entity, taskData map[string]interface{}) (task crm.Task, err error) {
	body, err := us.doPostEmptyHeaders(us.buildURL(crm.VersionUrl, "static/tasks"), taskData)
	if err != nil {
		return task, err
	}
	return task, json.Unmarshal(body, &task)
}

// GetField returns Field struct for a given type of entity & field
func (us *Uspacy) GetField(entityType crm.Entity, fieldType string) (field crm.Field, err error) {
	body, err := us.doGetEmptyHeaders(us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.FieldsUrl, entityType.GetUrl(), fieldType)))
	if err != nil {
		return field, err
	}
	return field, json.Unmarshal(body, &field)
}

// GetFields returns Fields struct for a given type of entity
func (us *Uspacy) GetFields(entityType crm.Entity) (fields crm.Fields, err error) {
	body, err := us.doGetEmptyHeaders(us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.FieldsUrl, entityType.GetUrl(), "")))
	if err != nil {
		return fields, err
	}
	return fields, json.Unmarshal(body, &fields)
}

// CreateFunnel returns created funnel
func (us *Uspacy) CreateFunnel(entityType crm.Entity, funnelData interface{}) (entityFunnel crm.Funnel, err error) {
	responseBody, err := us.doPostEmptyHeaders(us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.FunnelUrl, entityType.GetUrl())), funnelData)
	if err != nil {
		return entityFunnel, err
	}
	return entityFunnel, json.Unmarshal(responseBody, &entityFunnel)
}

// CreateFunnelStage returns created kanban stage
func (us *Uspacy) CreateFunnelStage(entityType crm.Entity, stageData interface{}) (kanbanStage crm.KanbanStage, err error) {
	responseBody, err := us.doPostEmptyHeaders(us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.KanbanStageUrl, entityType.GetUrl())), stageData)
	if err != nil {
		return kanbanStage, err
	}
	return kanbanStage, json.Unmarshal(responseBody, &kanbanStage)
}

// Move a funnel stage
func (us *Uspacy) MoveFunnelStage(entityType crm.Entity, entityId int64, stageId string) (err error) {
	_, err = us.doPostEmptyHeaders(us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.MoveKanbanStageUrl, entityType.GetUrl(), entityId, stageId)), nil)
	return err
}

// CreateCRMField in CRM entity returns created field
func (us *Uspacy) CreateCRMField(entityType crm.Entity, fieldData interface{}) (entityField crm.Field, err error) {
	responseBody, err := us.doPostEmptyHeaders(us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.FieldsUrl, entityType.GetUrl(), "")), fieldData)
	if err != nil {
		return entityField, err
	}
	return entityField, json.Unmarshal(responseBody, &entityField)
}

// CreateListValues returns arrey of values for given type of CRM list
func (us *Uspacy) CreateListValues(entityType crm.Entity, listName string, listValue interface{}) (lists []crm.List, err error) {
	responseBody, err := us.doPostEmptyHeaders(us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.ListsUrl, entityType.GetUrl(), listName)), listValue)
	if err != nil {
		return lists, err
	}
	return lists, json.Unmarshal(responseBody, &lists)
}
