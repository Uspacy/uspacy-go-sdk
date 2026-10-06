package api

import (
	"context"
	"fmt"

	"github.com/Uspacy/uspacy-go-sdk/v2/crm"
	"github.com/Uspacy/uspacy-go-sdk/v2/smartobjects"
)

// CreateSmartObject create smart object, retun created object and error
func (us *Uspacy) CreateSmartObject(ctx context.Context, fieldData smartobjects.SmartObjectCreateRequest, opts ...RequestOption) (smartobjects.CrmSmartObject, error) {
	responseBody, _, err := us.doPost(ctx, us.buildURL(crm.VersionUrl, crm.EntitiesUrl), fieldData, opts...)
	return decodeJSON[smartobjects.CrmSmartObject](responseBody, err)
}

// CreateSmartObjectEntity this method return any created object id, responce come and error
func (us *Uspacy) CreateSmartObjectEntity(ctx context.Context, tableName string, entityData map[string]any, opts ...RequestOption) (int64, int, error) {
	respBytes, code, err := us.doPost(ctx, us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.EntityUrl, tableName)), entityData, opts...)
	created, err := decodeJSON[createdID](respBytes, err)
	if err != nil {
		return 0, code, err
	}
	return created.ID, code, nil
}

// CreateSmartObjectField create field for selected smart object, retun created field and error
func (us *Uspacy) CreateSmartObjectField(ctx context.Context, tableName string, fieldData smartobjects.Field, opts ...RequestOption) (crm.Field, error) {
	responseBody, _, err := us.doPost(ctx, us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.CreateFieldUrl, tableName)), fieldData, opts...)
	return decodeJSON[crm.Field](responseBody, err)
}

// CreateSmartObjectListValues returns arrey of values for given type of CRM list
func (us *Uspacy) CreateSmartObjectListValues(ctx context.Context, tableName string, listName string, listValue any) ([]crm.List, error) {
	responseBody, _, err := us.doPost(ctx, us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.ListsUrl, tableName, listName)), listValue)
	return decodeJSON[[]crm.List](responseBody, err)
}

// CreateSmartObjectStage returns lwst of kanban stages
func (us *Uspacy) CreateSmartObjectStage(ctx context.Context, tableName string, stageData any, opts ...RequestOption) (crm.KanbanStage, error) {
	responseBody, _, err := us.doPost(ctx, us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.KanbanStageUrl, tableName, "")), stageData, opts...)
	return decodeJSON[crm.KanbanStage](responseBody, err)
}

// MoveSmartObjectFunnelStage moves a funnel stage
func (us *Uspacy) MoveSmartObjectFunnelStage(ctx context.Context, tableName string, entityId int64, stageId string, reason crm.KanbanFailReasonCRM, opts ...RequestOption) (err error) {
	_, _, err = us.doPost(ctx, us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.MoveKanbanStageUrl, tableName, entityId, stageId)), reason, opts...)
	return err
}

// GetSmartObjectFields returns Fields struct for a given table name of smart object
func (us *Uspacy) GetSmartObjectFields(ctx context.Context, tableName string) ([]crm.Field, error) {
	body, err := us.doGetEmptyHeaders(ctx, us.buildURL(crm.VersionUrl, fmt.Sprintf(smartobjects.FieldsUrl, tableName)))
	resp, err := decodeJSON[crm.Fields](body, err)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// GetSmartObjectStages list of smart object stages with given table name
func (us *Uspacy) GetSmartObjectStages(ctx context.Context, tableName string) ([]crm.KanbanStage, error) {
	body, err := us.doGetEmptyHeaders(ctx, us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.KanbanStageUrl, tableName, "")))
	resp, err := decodeJSON[crm.KanbanStages](body, err)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}
