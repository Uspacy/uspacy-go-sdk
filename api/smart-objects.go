package api

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Uspacy/uspacy-go-sdk/v2/crm"
	"github.com/Uspacy/uspacy-go-sdk/v2/smartobjects"
)

// CreateSmartObject create smart object, retun created object and error
func (us *Uspacy) CreateSmartObject(ctx context.Context, fieldData smartobjects.SmartObjectCreateRequest, opts ...RequestOption) (createdObject smartobjects.CrmSmartObject, err error) {
	responseBody, _, err := us.doPost(ctx, us.buildURL(crm.VersionUrl, crm.EntitiesUrl), fieldData, opts...)
	if err != nil {
		return createdObject, err
	}
	return createdObject, json.Unmarshal(responseBody, &createdObject)
}

// CreateSmartObjectEntity this method return any created object id, responce come and error
func (us *Uspacy) CreateSmartObjectEntity(ctx context.Context, tableName string, entityData map[string]any, opts ...RequestOption) (int64, int, error) {
	respBytes, code, err := us.doPost(ctx, us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.EntityUrl, tableName)), entityData, opts...)
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
func (us *Uspacy) CreateSmartObjectField(ctx context.Context, tableName string, fieldData smartobjects.Field, opts ...RequestOption) (entityField crm.Field, err error) {
	responseBody, _, err := us.doPost(ctx, us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.CreateFieldUrl, tableName)), fieldData, opts...)
	if err != nil {
		return entityField, err
	}
	return entityField, json.Unmarshal(responseBody, &entityField)
}

// CreateSmartObjectListValues returns arrey of values for given type of CRM list
func (us *Uspacy) CreateSmartObjectListValues(ctx context.Context, tableName string, listName string, listValue any) (lists []crm.List, err error) {
	responseBody, _, err := us.doPost(ctx, us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.ListsUrl, tableName, listName)), listValue)
	if err != nil {
		return lists, err
	}
	return lists, json.Unmarshal(responseBody, &lists)
}

// CreateSmartObjectStage returns lwst of kanban stages
func (us *Uspacy) CreateSmartObjectStage(ctx context.Context, tableName string, stageData any, opts ...RequestOption) (kanbanStage crm.KanbanStage, err error) {
	responseBody, _, err := us.doPost(ctx, us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.KanbanStageUrl, tableName, "")), stageData, opts...)
	if err != nil {
		return kanbanStage, err
	}
	return kanbanStage, json.Unmarshal(responseBody, &kanbanStage)
}

// MoveSmartObjectFunnelStage moves a funnel stage
func (us *Uspacy) MoveSmartObjectFunnelStage(ctx context.Context, tableName string, entityId int64, stageId string, reason crm.KanbanFailReasonCRM, opts ...RequestOption) (err error) {
	_, _, err = us.doPost(ctx, us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.MoveKanbanStageUrl, tableName, entityId, stageId)), reason, opts...)
	return err
}

// GetSmartObjectFields returns Fields struct for a given table name of smart object
func (us *Uspacy) GetSmartObjectFields(ctx context.Context, tableName string) (fields []crm.Field, err error) {
	body, err := us.doGetEmptyHeaders(ctx, us.buildURL(crm.VersionUrl, fmt.Sprintf(smartobjects.FieldsUrl, tableName)))
	if err != nil {
		return fields, err
	}
	var resp crm.Fields
	err = json.Unmarshal(body, &resp)
	return resp.Data, err
}

// GetSmartObjectStages list of smart object stages with given table name
func (us *Uspacy) GetSmartObjectStages(ctx context.Context, tableName string) (kanbanStages []crm.KanbanStage, err error) {
	body, err := us.doGetEmptyHeaders(ctx, us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.KanbanStageUrl, tableName, "")))
	if err != nil {
		return kanbanStages, err
	}
	var resp crm.KanbanStages
	err = json.Unmarshal(body, &resp)
	return resp.Data, err
}
