package api

import (
	"context"
	"fmt"

	"github.com/Uspacy/uspacy-go-sdk/v2/crm"
	"github.com/Uspacy/uspacy-go-sdk/v2/smartobjects"
)

// CreateSmartObject creates a smart object and returns it.
func (us *Uspacy) CreateSmartObject(ctx context.Context, fieldData smartobjects.SmartObjectCreateRequest, opts ...RequestOption) (smartobjects.CrmSmartObject, error) {
	responseBody, _, err := us.doPost(ctx, us.buildURL(crm.VersionUrl, crm.EntitiesUrl), fieldData, opts...)
	return decodeJSON[smartobjects.CrmSmartObject](responseBody, err)
}

// CreateSmartObjectEntity creates a smart object record and returns its ID and the HTTP status code.
func (us *Uspacy) CreateSmartObjectEntity(ctx context.Context, tableName string, entityData map[string]any, opts ...RequestOption) (int64, int, error) {
	respBytes, code, err := us.doPost(ctx, us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.EntityUrl, tableName)), entityData, opts...)
	created, err := decodeJSON[createdID](respBytes, err)
	if err != nil {
		return 0, code, err
	}
	return created.ID, code, nil
}

// CreateSmartObjectField creates a field of a smart object and returns it.
func (us *Uspacy) CreateSmartObjectField(ctx context.Context, tableName string, fieldData smartobjects.Field, opts ...RequestOption) (crm.Field, error) {
	responseBody, _, err := us.doPost(ctx, us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.CreateFieldUrl, tableName)), fieldData, opts...)
	return decodeJSON[crm.Field](responseBody, err)
}

// CreateSmartObjectListValues adds values to a smart object list field and returns the field's values.
func (us *Uspacy) CreateSmartObjectListValues(ctx context.Context, tableName string, listName string, listValue any, opts ...RequestOption) ([]crm.List, error) {
	responseBody, _, err := us.doPost(ctx, us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.ListsUrl, tableName, listName)), listValue, opts...)
	return decodeJSON[[]crm.List](responseBody, err)
}

// CreateSmartObjectStage creates a kanban stage of a smart object and returns it.
func (us *Uspacy) CreateSmartObjectStage(ctx context.Context, tableName string, stageData any, opts ...RequestOption) (crm.KanbanStage, error) {
	responseBody, _, err := us.doPost(ctx, us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.KanbanStageUrl, tableName, "")), stageData, opts...)
	return decodeJSON[crm.KanbanStage](responseBody, err)
}

// MoveSmartObjectFunnelStage moves a funnel stage
func (us *Uspacy) MoveSmartObjectFunnelStage(ctx context.Context, tableName string, entityId int64, stageId string, reason crm.KanbanFailReasonCRM, opts ...RequestOption) (err error) {
	_, _, err = us.doPost(ctx, us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.MoveKanbanStageUrl, tableName, entityId, stageId)), reason, opts...)
	return err
}

// GetSmartObjectFields returns the fields of a smart object.
func (us *Uspacy) GetSmartObjectFields(ctx context.Context, tableName string, opts ...RequestOption) ([]crm.Field, error) {
	body, err := us.doGet(ctx, us.buildURL(crm.VersionUrl, fmt.Sprintf(smartobjects.FieldsUrl, tableName)), opts...)
	resp, err := decodeJSON[crm.Fields](body, err)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// GetSmartObjectStages returns the kanban stages of a smart object.
func (us *Uspacy) GetSmartObjectStages(ctx context.Context, tableName string, opts ...RequestOption) ([]crm.KanbanStage, error) {
	body, err := us.doGet(ctx, us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.KanbanStageUrl, tableName, "")), opts...)
	resp, err := decodeJSON[crm.KanbanStages](body, err)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}
