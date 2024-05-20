package api

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/Uspacy/uspacy-go-sdk/crm"
)

// CreateEntity this method does not return any object, just error
func (us *Uspacy) CreateEntity(entityType crm.Entity, entityData map[string]interface{}) (int64, int, error) {
	respBytes, code, err := us.doPostEmptyHeaders(us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.EntityUrl, entityType.GetUrl())), entityData)
	if err != nil {
		return 0, code, err
	}

	var respData struct {
		ID int64 `json:"id"`
	}
	if err := json.Unmarshal(respBytes, &respData); err != nil {
		return 0, code, err
	}

	return respData.ID, code, nil
}

// GetEntities this method return arrey of entities present in crm and error
func (us *Uspacy) GetCrmEntitiesList() (entities []crm.CrmEntities, err error) {
	body, err := us.doGetEmptyHeaders(us.buildURL(crm.VersionUrl, crm.EntitiesUrl))
	if err != nil {
		return entities, err
	}
	var resp = crm.CrmEntitiesList{}
	err = json.Unmarshal(body, &resp)
	return resp.Data, err
}

// GetEntities this method return arrey of entities present in crm and error
func (us *Uspacy) GetEntities(entityType crm.Entity, params url.Values) (entities crm.CRMEntity, err error) {
	body, err := us.doGetEmptyHeaders(us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.EntityUrl, entityType.GetUrl())) + "?" + params.Encode())
	if err != nil {
		return entities, err
	}
	return entities, json.Unmarshal(body, &entities)
}

// GetEntities this method return arrey of objects wifh all fields and error
func (us *Uspacy) GetCRMEntitiesForExport(entityType crm.Entity, params url.Values) (entities crm.CRMEntityForExport, err error) {
	var entityRoute string
	switch entityType {
	case crm.TasksNum:
		entityRoute = fmt.Sprintf(crm.TaskUrl, "")
	case crm.ProductsNum:
		entityRoute = fmt.Sprintf(crm.ProductsUrl, "")
	default:
		entityRoute = fmt.Sprintf(crm.EntityUrl, entityType.GetUrl())
	}
	body, err := us.doGetEmptyHeaders(us.buildURL(crm.VersionUrl, entityRoute) + "?" + params.Encode())
	if err != nil {
		return entities, err
	}
	return entities, json.Unmarshal(body, &entities)
}

// GetEntities this method return arrey of objects and error
func (us *Uspacy) GetContacts(entityType crm.Entity, params url.Values) (entities crm.Contacts, err error) {
	body, err := us.doGetEmptyHeaders(us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.EntityUrl, entityType.GetUrl())) + "?" + params.Encode())
	if err != nil {
		return entities, err
	}
	return entities, json.Unmarshal(body, &entities)
}

// PatchEntity this method does not return any object, just error
func (us *Uspacy) PatchEntity(entityType string, id string, entityData map[string]interface{}) error {
	_, err := us.doPatchEmptyHeaders(us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.EntityUrl, entityType))+id, entityData)
	if err != nil {
		return err
	}
	return nil
}

// CreateContact returns created contact object
func (us *Uspacy) CreateContact(entityType crm.Entity, contactData map[string]interface{}) (contact crm.Contact, err error) {
	body, _, err := us.doPostEmptyHeaders(us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.EntityUrl, entityType.GetUrl())), contactData)
	if err != nil {
		return contact, err
	}
	return contact, json.Unmarshal(body, &contact)
}

// CreateCompany returns created company object
func (us *Uspacy) CreateCompany(entityType crm.Entity, companyData map[string]interface{}) (company crm.Company, err error) {
	body, _, err := us.doPostEmptyHeaders(us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.EntityUrl, entityType.GetUrl())), companyData)
	if err != nil {
		return company, err
	}
	return company, json.Unmarshal(body, &company)
}

// CreateLeads returns created lead object
func (us *Uspacy) CreateLead(entityType crm.Entity, leadData map[string]interface{}) (lead crm.Lead, err error) {
	body, _, err := us.doPostEmptyHeaders(us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.EntityUrl, entityType.GetUrl())), leadData)
	if err != nil {
		return lead, err
	}
	return lead, json.Unmarshal(body, &lead)
}

// CreateDeals returns created deal object
func (us *Uspacy) CreateDeal(entityType crm.Entity, dealData map[string]interface{}) (deal crm.Deal, err error) {
	body, _, err := us.doPostEmptyHeaders(us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.EntityUrl, entityType.GetUrl())), dealData)
	if err != nil {
		return deal, err
	}
	return deal, json.Unmarshal(body, &deal)
}

// CreateTask returns created task object
func (us *Uspacy) CreateTaskCRM(entityType crm.Entity, taskData map[string]interface{}) (task crm.Task, err error) {
	body, _, err := us.doPostEmptyHeaders(us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.TaskUrl, "")), taskData)
	if err != nil {
		return task, err
	}
	return task, json.Unmarshal(body, &task)
}

// PatchTaskCrm returns created task object
func (us *Uspacy) PatchTaskCrm(entityType crm.Entity, id string, taskData map[string]interface{}) (task crm.Task, err error) {
	body, err := us.doPatchEmptyHeaders(us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.TaskUrl, id)), taskData)
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

// DeleteField delete selected field for given type of entity
func (us *Uspacy) DeleteField(entityType crm.Entity, codeField string) (err error) {
	_, err = us.doDeleteEmptyHeaders(us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.FieldsUrl, entityType.GetUrl(), codeField)))
	return err
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
	responseBody, _, err := us.doPostEmptyHeaders(us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.FunnelUrl, entityType.GetUrl())), funnelData)
	if err != nil {
		return entityFunnel, err
	}
	return entityFunnel, json.Unmarshal(responseBody, &entityFunnel)
}

// GetFunnels returns funnels by entityType
func (us *Uspacy) GetFunnels(entityType crm.Entity) (funnels crm.FunnelsById, err error) {
	responseBody, err := us.doGetEmptyHeaders(us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.FunnelUrl, entityType.GetUrl())))
	if err != nil {
		return funnels, err
	}
	return funnels, json.Unmarshal(responseBody, &funnels)
}

// CreateFunnelStage returns created kanban stage
func (us *Uspacy) CreateFunnelStage(entityType crm.Entity, stageData interface{}) (kanbanStage crm.KanbanStage, err error) {
	responseBody, _, err := us.doPostEmptyHeaders(us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.KanbanStageUrl, entityType.GetUrl(), "")), stageData)
	if err != nil {
		return kanbanStage, err
	}
	return kanbanStage, json.Unmarshal(responseBody, &kanbanStage)
}

// GetFunnelStage returns all kanban stages
func (us *Uspacy) GetAllFunnelStages(entityType crm.Entity) (kanbanStages crm.KanbanStages, err error) {
	responseBody, err := us.doGetEmptyHeaders(us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.KanbanStageUrl, entityType.GetUrl(), "")))
	if err != nil {
		return kanbanStages, err
	}
	return kanbanStages, json.Unmarshal(responseBody, &kanbanStages)
}

// GetFunnelStage returns kanban stage
func (us *Uspacy) GetFunnelStageDyId(entityType crm.Entity, id int) (kanbanStages crm.KanbanStages, err error) {
	responseBody, err := us.doGetEmptyHeaders(us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.KanbanStageUrl, entityType.GetUrl(), fmt.Sprintf(crm.StageByFunnelIdUrl, id))))
	if err != nil {
		return kanbanStages, err
	}
	return kanbanStages, json.Unmarshal(responseBody, &kanbanStages)
}

// PatchFunnelStage returns kanban stage
func (us *Uspacy) PatchFunnelStage(entityType crm.Entity, id int, stage crm.FunnelStage) (kanbanStage crm.KanbanStage, err error) {
	responseBody, err := us.doPatchEmptyHeaders(us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.KanbanStageUrl, entityType.GetUrl(), id)), stage)
	if err != nil {
		return kanbanStage, err
	}
	return kanbanStage, json.Unmarshal(responseBody, &kanbanStage)
}

// Move a funnel stage
func (us *Uspacy) MoveFunnelStage(entityType crm.Entity, entityId int64, stageId string, reason crm.KanbanFailReasonCRM) (err error) {
	_, _, err = us.doPostEmptyHeaders(us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.MoveKanbanStageUrl, entityType.GetUrl(), entityId, stageId)), reason)
	return err
}

// CreateCRMField in CRM entity returns created field
func (us *Uspacy) CreateCRMField(entityType crm.Entity, fieldData interface{}) (entityField crm.Field, err error) {
	responseBody, _, err := us.doPostEmptyHeaders(us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.CreateFieldUrl, entityType.GetUrl())), fieldData)
	if err != nil {
		return entityField, err
	}
	return entityField, json.Unmarshal(responseBody, &entityField)
}

// GetListValues returns arrey of values for given type of CRM list
func (us *Uspacy) GetListValues(entityType crm.Entity, listName string) (lists []crm.List, err error) {
	responseBody, err := us.doGetEmptyHeaders(us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.ListsUrl, entityType.GetUrl(), listName)))
	if err != nil {
		return lists, err
	}
	return lists, json.Unmarshal(responseBody, &lists)
}

// CreateListValues returns arrey of values for given type of CRM list
func (us *Uspacy) CreateListValues(entityType crm.Entity, listName string, listValue interface{}) (lists []crm.List, err error) {
	responseBody, _, err := us.doPostEmptyHeaders(us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.ListsUrl, entityType.GetUrl(), listName)), listValue)
	if err != nil {
		return lists, err
	}
	return lists, json.Unmarshal(responseBody, &lists)
}

// CreateFailReasons returns all reasons for funnel with failWrite.ID
func (us *Uspacy) CreateFailReasons(failReason crm.Reason) (reasons crm.Reason, err error) {
	responseBody, _, err := us.doPostEmptyHeaders(us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.ReasonsUrl, failReason.ID)), crm.FailWrite{
		Title: failReason.Title,
		Sort:  failReason.Sort,
		Type:  "FAIL",
	})
	if err != nil {
		return reasons, err
	}
	return reasons, json.Unmarshal(responseBody, &reasons)
}

// CreateCall returns created call
func (us *Uspacy) CreateCall(callValue crm.Call) (call crm.Call, err error) {
	responseBody, _, err := us.doPostEmptyHeaders(us.buildURL(crm.VersionUrl, crm.CallUrl), callValue)
	if err != nil {
		return call, err
	}
	return call, json.Unmarshal(responseBody, &call)
}
