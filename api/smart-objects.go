package api

import (
	"encoding/json"
	"fmt"

	"github.com/Uspacy/uspacy-go-sdk/crm"
	"github.com/Uspacy/uspacy-go-sdk/smartobjects"
)

// CreateSmartObject create smart object, retun created object and error
func (us *Uspacy) CreateSmartObject(fieldData smartobjects.SmartObjectCreateRequest, headers ...map[string]string) (createdObject smartobjects.CrmSmartObject, err error) {
	responseBody, _, err := us.doPost(us.buildURL(crm.VersionUrl, crm.EntitiesUrl), fieldData, headers...)
	if err != nil {
		return createdObject, err
	}
	return createdObject, json.Unmarshal(responseBody, &createdObject)
}

// CreateSmartObjectEntity this method return any created object id, responce come and error
func (us *Uspacy) CreateSmartObjectEntity(tableName string, entityData map[string]interface{}, headers ...map[string]string) (int64, int, error) {
	respBytes, code, err := us.doPost(us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.EntityUrl, tableName)), entityData, headers...)
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

// CreateSmartObjectField create field for selected smart object, retun created field and error
func (us *Uspacy) CreateSmartObjectField(tableName string, fieldData smartobjects.Field, headers ...map[string]string) (entityField crm.Field, err error) {
	responseBody, _, err := us.doPost(us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.CreateFieldUrl, tableName)), fieldData, headers...)
	if err != nil {
		return entityField, err
	}
	return entityField, json.Unmarshal(responseBody, &entityField)
}

// CreateListValues returns arrey of values for given type of CRM list
func (us *Uspacy) CreateSmartObjectListValues(tableName string, listName string, listValue interface{}) (lists []crm.List, err error) {
	responseBody, _, err := us.doPost(us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.ListsUrl, tableName, listName)), listValue)
	if err != nil {
		return lists, err
	}
	return lists, json.Unmarshal(responseBody, &lists)
}

// CreateSmartObjectStage returns lwst of kanban stages
func (us *Uspacy) CreateSmartObjectStage(tableName string, stageData interface{}, headers ...map[string]string) (kanbanStage crm.KanbanStage, err error) {
	responseBody, _, err := us.doPost(us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.KanbanStageUrl, tableName, "")), stageData, headers...)
	if err != nil {
		return kanbanStage, err
	}
	return kanbanStage, json.Unmarshal(responseBody, &kanbanStage)
}

// Move a funnel stage
func (us *Uspacy) MoveSmartObjectFunnelStage(tableName string, entityId int64, stageId string, reason crm.KanbanFailReasonCRM, headers ...map[string]string) (err error) {
	_, _, err = us.doPost(us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.MoveKanbanStageUrl, tableName, entityId, stageId)), reason, headers...)
	return err
}

// GetSmartObjectFields returns Fields struct for a given table name of smart object
func (us *Uspacy) GetSmartObjectFields(tableName string) (fields []crm.Field, err error) {
	body, err := us.doGetEmptyHeaders(us.buildURL(crm.VersionUrl, fmt.Sprintf(smartobjects.FieldsUrl, tableName)))
	if err != nil {
		return fields, err
	}
	var resp crm.Fields
	err = json.Unmarshal(body, &resp)
	return resp.Data, err
}

// GetSmartObjectStages list of smart object stages with given table name
func (us *Uspacy) GetSmartObjectStages(tableName string) (kanbanStages []crm.KanbanStage, err error) {
	body, err := us.doGetEmptyHeaders(us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.KanbanStageUrl, tableName, "")))
	if err != nil {
		return kanbanStages, err
	}
	var resp crm.KanbanStages
	err = json.Unmarshal(body, &resp)
	return resp.Data, err
}
