package api

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/Uspacy/uspacy-go-sdk/crm"
)

// CreateEntity this method does not return any object, just error
func (us *Uspacy) CreateEntity(entityType string, entityData map[string]any, headers ...map[string]string) (int64, int, error) {
	respBytes, code, err := us.doPost(us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.EntityUrl, entityType)), entityData, headers...)
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
func (us *Uspacy) GetEntities(entityType string, params url.Values) (entities crm.CRMEntity, err error) {
	body, err := us.doGetEmptyHeaders(us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.EntityUrl, entityType)) + "?" + params.Encode())
	if err != nil {
		return entities, err
	}
	return entities, json.Unmarshal(body, &entities)
}

// GetEntities this method return arrey of objects wifh all fields and error
func (us *Uspacy) GetCRMEntitiesForExport(entityType string, params url.Values) (entities crm.CRMEntityForExport, err error) {
	var entityRoute string
	switch entityType {
	case crm.TasksNum.GetUrl():
		entityRoute = fmt.Sprintf(crm.TaskUrl, "")
	case crm.ProductsNum.GetUrl():
		entityRoute = fmt.Sprintf(crm.ProductsUrl, "")
	default:
		entityRoute = fmt.Sprintf(crm.EntityUrl, entityType)
	}
	body, err := us.doGetEmptyHeaders(us.buildURL(crm.VersionUrl, entityRoute) + "?" + params.Encode())
	if err != nil {
		return entities, err
	}
	return entities, json.Unmarshal(body, &entities)
}

// GetEntities this method return arrey of objects and error
func (us *Uspacy) GetContacts(params url.Values) (entities crm.Contacts, err error) {
	body, err := us.doGetEmptyHeaders(us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.EntityUrl, crm.ContactsNum.GetUrl())) + "?" + params.Encode())
	if err != nil {
		return entities, err
	}
	return entities, json.Unmarshal(body, &entities)
}

// PatchEntity this method does not return any object, just error
func (us *Uspacy) PatchEntity(entityType string, id string, entityData map[string]any) error {
	_, err := us.doPatchEmptyHeaders(us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.EntityUrl, entityType))+id, entityData)
	if err != nil {
		return err
	}
	return nil
}

// PatchEntity this method does not return any object, just error
func (us *Uspacy) EntityMassEdit(entityType string, entityData map[string]any) error {
	_, err := us.doPatchEmptyHeaders(us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.EntityUrl, entityType))+"mass_edit", entityData)
	if err != nil {
		return err
	}
	return nil
}

// CreateContact returns created contact object
func (us *Uspacy) CreateContact(contactData map[string]any, headers ...map[string]string) (contact crm.Contact, err error) {
	body, _, err := us.doPost(us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.EntityUrl, crm.ContactsNum.GetUrl())), contactData, headers...)
	if err != nil {
		return contact, err
	}
	return contact, json.Unmarshal(body, &contact)
}

// CreateCompany returns created company object
func (us *Uspacy) CreateCompany(companyData map[string]any, headers ...map[string]string) (company crm.Company, err error) {
	body, _, err := us.doPost(us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.EntityUrl, crm.CompaniesNum.GetUrl())), companyData, headers...)
	if err != nil {
		return company, err
	}
	return company, json.Unmarshal(body, &company)
}

// CreateLeads returns created lead object
func (us *Uspacy) CreateLead(leadData map[string]any, headers ...map[string]string) (lead crm.Lead, err error) {
	body, _, err := us.doPost(us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.EntityUrl, crm.LeadsNum.GetUrl())), leadData, headers...)
	if err != nil {
		return lead, err
	}
	return lead, json.Unmarshal(body, &lead)
}

// CreateDeals returns created deal object
func (us *Uspacy) CreateDeal(dealData map[string]any, headers ...map[string]string) (deal crm.Deal, err error) {
	body, _, err := us.doPost(us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.EntityUrl, crm.DealsNum.GetUrl())), dealData, headers...)
	if err != nil {
		return deal, err
	}
	return deal, json.Unmarshal(body, &deal)
}

// CreateTask returns created task object
func (us *Uspacy) CreateTaskCRM(taskData map[string]any, headers ...map[string]string) (task crm.Task, err error) {
	body, _, err := us.doPost(us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.TaskUrl, "")), taskData, headers...)
	if err != nil {
		return task, err
	}
	return task, json.Unmarshal(body, &task)
}

// PatchTaskCrm returns created task object
func (us *Uspacy) PatchTaskCrm(id string, taskData map[string]any) (task crm.Task, err error) {
	body, err := us.doPatchEmptyHeaders(us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.TaskUrl, id)), taskData)
	if err != nil {
		return task, err
	}
	return task, json.Unmarshal(body, &task)
}

// GetField returns Field struct for a given type of entity & field
func (us *Uspacy) GetField(entityType string, fieldType string) (field crm.Field, err error) {
	body, err := us.doGetEmptyHeaders(us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.FieldsUrl, entityType, fieldType)))
	if err != nil {
		return field, err
	}
	return field, json.Unmarshal(body, &field)
}

// DeleteField delete selected field for given type of entity
func (us *Uspacy) DeleteField(entityType string, codeField string) (err error) {
	_, err = us.doDeleteEmptyHeaders(us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.FieldsUrl, entityType, codeField)), nil)
	return err
}

// GetFields returns Fields struct for a given type of entity
func (us *Uspacy) GetFields(entityType string) (fields []crm.Field, err error) {
	body, err := us.doGetEmptyHeaders(us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.FieldsUrl, entityType, "")))
	if err != nil {
		return fields, err
	}
	var resp crm.Fields
	return resp.Data, json.Unmarshal(body, &resp)
}

// CreateFunnel returns created funnel
func (us *Uspacy) CreateFunnel(entityType string, funnelData any, headers ...map[string]string) (entityFunnel crm.Funnel, err error) {
	responseBody, _, err := us.doPost(us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.FunnelUrl, entityType)), funnelData, headers...)
	if err != nil {
		return entityFunnel, err
	}
	return entityFunnel, json.Unmarshal(responseBody, &entityFunnel)
}

// GetFunnels returns funnels by entityType
func (us *Uspacy) GetFunnels(entityType string) (funnels crm.FunnelsById, err error) {
	responseBody, err := us.doGetEmptyHeaders(us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.FunnelUrl, entityType)))
	if err != nil {
		return funnels, err
	}
	return funnels, json.Unmarshal(responseBody, &funnels)
}

// CreateFunnelStage returns created kanban stage
func (us *Uspacy) CreateFunnelStage(entityType string, stageData any, headers ...map[string]string) (kanbanStage crm.KanbanStage, err error) {
	responseBody, _, err := us.doPost(us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.KanbanStageUrl, entityType, "")), stageData, headers...)
	if err != nil {
		return kanbanStage, err
	}
	return kanbanStage, json.Unmarshal(responseBody, &kanbanStage)
}

// GetFunnelStage returns all kanban stages
func (us *Uspacy) GetAllFunnelStages(entityType string) (kanbanStages []crm.KanbanStage, err error) {
	responseBody, err := us.doGetEmptyHeaders(us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.KanbanStageUrl, entityType, "")))
	if err != nil {
		return kanbanStages, err
	}
	var resp crm.KanbanStages
	return resp.Data, json.Unmarshal(responseBody, &resp)
}

// GetFunnelStage returns kanban stage
func (us *Uspacy) GetFunnelStageDyId(entityType string, id int) (kanbanStages crm.KanbanStages, err error) {
	responseBody, err := us.doGetEmptyHeaders(us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.KanbanStageUrl, entityType, fmt.Sprintf(crm.StageByFunnelIdUrl, id))))
	if err != nil {
		return kanbanStages, err
	}
	return kanbanStages, json.Unmarshal(responseBody, &kanbanStages)
}

// PatchFunnelStage returns kanban stage
func (us *Uspacy) PatchFunnelStage(entityType string, id int, stage crm.FunnelStage) (kanbanStage crm.KanbanStage, err error) {
	responseBody, err := us.doPatchEmptyHeaders(us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.KanbanStageUrl, entityType, id)), stage)
	if err != nil {
		return kanbanStage, err
	}
	return kanbanStage, json.Unmarshal(responseBody, &kanbanStage)
}

// Move a funnel stage
func (us *Uspacy) MoveFunnelStage(entityType string, entityId int64, stageId string, reason crm.KanbanFailReasonCRM, headers ...map[string]string) (err error) {
	_, _, err = us.doPost(us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.MoveKanbanStageUrl, entityType, entityId, stageId)), reason, headers...)
	return err
}

// CreateCRMField in CRM entity returns created field
func (us *Uspacy) CreateCRMField(entityType string, fieldData any, headers ...map[string]string) (entityField crm.Field, err error) {
	responseBody, _, err := us.doPost(us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.CreateFieldUrl, entityType)), fieldData, headers...)
	if err != nil {
		return entityField, err
	}
	return entityField, json.Unmarshal(responseBody, &entityField)
}

// GetListValues returns arrey of values for given type of CRM list
func (us *Uspacy) GetListValues(entityType, listName string) (lists []crm.List, err error) {
	responseBody, err := us.doGetEmptyHeaders(us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.ListsUrl, entityType, listName)))
	if err != nil {
		return lists, err
	}
	return lists, json.Unmarshal(responseBody, &lists)
}

// CreateListValues returns arrey of values for given type of CRM list
func (us *Uspacy) CreateListValues(entityType, listName string, listValue any, headers ...map[string]string) (lists []crm.List, err error) {
	responseBody, _, err := us.doPost(us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.ListsUrl, entityType, listName)), listValue, headers...)
	if err != nil {
		return lists, err
	}
	return lists, json.Unmarshal(responseBody, &lists)
}

// CreateFailReasons returns all reasons for funnel with failWrite.ID
func (us *Uspacy) CreateFailReasons(failReason crm.Reason, headers ...map[string]string) (reasons crm.Reason, err error) {
	responseBody, _, err := us.doPost(us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.ReasonsUrl, failReason.ID)), crm.FailWrite{
		Title: failReason.Title,
		Sort:  failReason.Sort,
		Type:  "FAIL",
	}, headers...)
	if err != nil {
		return reasons, err
	}
	return reasons, json.Unmarshal(responseBody, &reasons)
}

// CreateCall returns created call
func (us *Uspacy) CreateCall(callValue crm.Call, headers ...map[string]string) (call crm.Call, err error) {
	responseBody, _, err := us.doPost(us.buildURL(crm.VersionUrl, crm.CallUrl), callValue, headers...)
	if err != nil {
		return call, err
	}
	return call, json.Unmarshal(responseBody, &call)
}
