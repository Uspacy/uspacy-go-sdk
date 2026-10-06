package api

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/Uspacy/uspacy-go-sdk/v2/crm"
)

// CreateEntity creates a CRM record and returns its ID and the HTTP status code.
func (us *Uspacy) CreateEntity(ctx context.Context, entityType string, entityData map[string]any, opts ...RequestOption) (int64, int, error) {
	respBytes, code, err := us.doPost(ctx, us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.EntityUrl, entityType)), entityData, opts...)
	created, err := decodeJSON[createdID](respBytes, err)
	if err != nil {
		return 0, code, err
	}
	return created.ID, code, nil
}

// GetCrmEntitiesList returns the entity types available in the CRM.
func (us *Uspacy) GetCrmEntitiesList(ctx context.Context, opts ...RequestOption) ([]crm.CrmEntities, error) {
	body, err := us.doGet(ctx, us.buildURL(crm.VersionUrl, crm.EntitiesUrl), opts...)
	resp, err := decodeJSON[crm.CrmEntitiesList](body, err)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// GetEntities returns the records of a CRM entity type matching params.
func (us *Uspacy) GetEntities(ctx context.Context, entityType string, params url.Values, opts ...RequestOption) (crm.CRMEntity, error) {
	body, err := us.doGet(ctx, us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.EntityUrl, entityType))+"?"+params.Encode(), opts...)
	return decodeJSON[crm.CRMEntity](body, err)
}

// GetCRMEntitiesForExport returns the records of a CRM entity type with all their fields.
func (us *Uspacy) GetCRMEntitiesForExport(ctx context.Context, entityType string, params url.Values, opts ...RequestOption) (crm.CRMEntityForExport, error) {
	var entityRoute string
	switch entityType {
	case crm.ProductsNum.GetUrl():
		entityRoute = fmt.Sprintf(crm.ProductsUrl, "")
	default:
		entityRoute = fmt.Sprintf(crm.EntityUrl, entityType)
	}
	body, err := us.doGet(ctx, us.buildURL(crm.VersionUrl, entityRoute)+"?"+params.Encode(), opts...)
	return decodeJSON[crm.CRMEntityForExport](body, err)
}

// GetContacts returns an array of contact objects and an error.
func (us *Uspacy) GetContacts(ctx context.Context, params url.Values, opts ...RequestOption) (crm.Contacts, error) {
	body, err := us.doGet(ctx, us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.EntityUrl, crm.ContactsNum.GetUrl()))+"?"+params.Encode(), opts...)
	return decodeJSON[crm.Contacts](body, err)
}

// GetDeals returns an array of deal objects and an error.
func (us *Uspacy) GetDeals(ctx context.Context, params url.Values, opts ...RequestOption) (crm.Deals, error) {
	body, err := us.doGet(ctx, us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.EntityUrl, crm.DealsNum.GetUrl()))+"?"+params.Encode(), opts...)
	return decodeJSON[crm.Deals](body, err)
}

// GetLeads returns an array of lead objects and an error.
func (us *Uspacy) GetLeads(ctx context.Context, params url.Values, opts ...RequestOption) (crm.Leads, error) {
	body, err := us.doGet(ctx, us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.EntityUrl, crm.LeadsNum.GetUrl()))+"?"+params.Encode(), opts...)
	return decodeJSON[crm.Leads](body, err)
}

// GetList returns raw response for CRM entities with filters
// This method is useful for searching entities with custom filters
// entityType should be one of: crm.LeadsNum.GetUrl(), crm.DealsNum.GetUrl(), crm.ContactsNum.GetUrl(), crm.CompaniesNum.GetUrl()
func (us *Uspacy) GetList(ctx context.Context, entityType string, params url.Values, opts ...RequestOption) ([]byte, error) {
	url := us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.EntityUrl, entityType))
	if len(params) > 0 {
		url += "?" + params.Encode()
	}

	return us.doGet(ctx, url, opts...)
}

// GetEntity returns one record as raw JSON: GET crm/v1/entities/{entityType}/{id}.
// The id is its own path segment so the leading slash is preserved.
// A non-2xx response is returned as *HTTPError.
func (us *Uspacy) GetEntity(ctx context.Context, entityType string, id int64, opts ...RequestOption) ([]byte, error) {
	return us.doGet(ctx, us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.EntityUrl, entityType), strconv.FormatInt(id, 10)), opts...)
}

// PatchEntity updates a CRM record.
func (us *Uspacy) PatchEntity(ctx context.Context, entityType string, id string, entityData map[string]any, opts ...RequestOption) error {
	_, err := us.doPatch(ctx, us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.EntityUrl, entityType), id), entityData, opts...)
	if err != nil {
		return err
	}
	return nil
}

// EntityMassEdit updates several CRM records of one entity type at once.
func (us *Uspacy) EntityMassEdit(ctx context.Context, entityType string, entityData map[string]any, opts ...RequestOption) error {
	_, err := us.doPatch(ctx, us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.EntityUrl, entityType), "mass_edit"), entityData, opts...)
	if err != nil {
		return err
	}
	return nil
}

// CreateContact returns created contact object
func (us *Uspacy) CreateContact(ctx context.Context, contactData map[string]any, opts ...RequestOption) (crm.Contact, error) {
	body, _, err := us.doPost(ctx, us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.EntityUrl, crm.ContactsNum.GetUrl())), contactData, opts...)
	return decodeJSON[crm.Contact](body, err)
}

// CreateCompany returns created company object
func (us *Uspacy) CreateCompany(ctx context.Context, companyData map[string]any, opts ...RequestOption) (crm.Company, error) {
	body, _, err := us.doPost(ctx, us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.EntityUrl, crm.CompaniesNum.GetUrl())), companyData, opts...)
	return decodeJSON[crm.Company](body, err)
}

// CreateLead returns the created lead.
func (us *Uspacy) CreateLead(ctx context.Context, leadData map[string]any, opts ...RequestOption) (crm.Lead, error) {
	body, _, err := us.doPost(ctx, us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.EntityUrl, crm.LeadsNum.GetUrl())), leadData, opts...)
	return decodeJSON[crm.Lead](body, err)
}

// CreateDeal returns the created deal.
func (us *Uspacy) CreateDeal(ctx context.Context, dealData map[string]any, opts ...RequestOption) (crm.Deal, error) {
	body, _, err := us.doPost(ctx, us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.EntityUrl, crm.DealsNum.GetUrl())), dealData, opts...)
	return decodeJSON[crm.Deal](body, err)
}

// GetField returns Field struct for a given type of entity & field
func (us *Uspacy) GetField(ctx context.Context, entityType string, fieldType string, opts ...RequestOption) (crm.Field, error) {
	body, err := us.doGet(ctx, us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.FieldsUrl, entityType, fieldType)), opts...)
	return decodeJSON[crm.Field](body, err)
}

// DeleteField deletes a field of a CRM entity type.
func (us *Uspacy) DeleteField(ctx context.Context, entityType string, codeField string, opts ...RequestOption) (err error) {
	_, err = us.doDelete(ctx, us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.FieldsUrl, entityType, codeField)), nil, opts...)
	return err
}

// GetFields returns the fields of an entity type: GET crm/v1/entities/{entityType}/fields.
// A non-2xx response is returned as *HTTPError.
func (us *Uspacy) GetFields(ctx context.Context, entityType string, opts ...RequestOption) ([]crm.Field, error) {
	body, err := us.doGet(ctx, us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.FieldsUrl, entityType, "")), opts...)
	resp, err := decodeJSON[crm.Fields](body, err)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// CreateFunnel returns created funnel
func (us *Uspacy) CreateFunnel(ctx context.Context, entityType string, funnelData any, opts ...RequestOption) (crm.Funnel, error) {
	responseBody, _, err := us.doPost(ctx, us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.FunnelUrl, entityType)), funnelData, opts...)
	return decodeJSON[crm.Funnel](responseBody, err)
}

// GetFunnels returns funnels by entityType
func (us *Uspacy) GetFunnels(ctx context.Context, entityType string, opts ...RequestOption) (crm.FunnelsById, error) {
	responseBody, err := us.doGet(ctx, us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.FunnelUrl, entityType)), opts...)
	return decodeJSON[crm.FunnelsById](responseBody, err)
}

// CreateFunnelStage returns created kanban stage
func (us *Uspacy) CreateFunnelStage(ctx context.Context, entityType string, stageData any, opts ...RequestOption) (crm.KanbanStage, error) {
	responseBody, _, err := us.doPost(ctx, us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.KanbanStageUrl, entityType, "")), stageData, opts...)
	return decodeJSON[crm.KanbanStage](responseBody, err)
}

// GetAllFunnelStages returns all kanban stages
func (us *Uspacy) GetAllFunnelStages(ctx context.Context, entityType string, opts ...RequestOption) ([]crm.KanbanStage, error) {
	responseBody, err := us.doGet(ctx, us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.KanbanStageUrl, entityType, "")), opts...)
	resp, err := decodeJSON[crm.KanbanStages](responseBody, err)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// GetFunnelStageById returns kanban stage by id.
//
// Renamed in v2 from GetFunnelStageDyId.
func (us *Uspacy) GetFunnelStageById(ctx context.Context, entityType string, id int, opts ...RequestOption) (crm.KanbanStages, error) {
	responseBody, err := us.doGet(ctx, us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.KanbanStageUrl, entityType, fmt.Sprintf(crm.StageByFunnelIdUrl, id))), opts...)
	return decodeJSON[crm.KanbanStages](responseBody, err)
}

// PatchFunnelStage returns kanban stage
func (us *Uspacy) PatchFunnelStage(ctx context.Context, entityType string, id int, stage crm.FunnelStage, opts ...RequestOption) (crm.KanbanStage, error) {
	responseBody, err := us.doPatch(ctx, us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.KanbanStageUrl, entityType, id)), stage, opts...)
	return decodeJSON[crm.KanbanStage](responseBody, err)
}

// MoveFunnelStage moves a funnel stage
func (us *Uspacy) MoveFunnelStage(ctx context.Context, entityType string, entityId int64, stageId string, reason crm.KanbanFailReasonCRM, opts ...RequestOption) (err error) {
	_, _, err = us.doPost(ctx, us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.MoveKanbanStageUrl, entityType, entityId, stageId)), reason, opts...)
	return err
}

// CreateCRMField in CRM entity returns created field
func (us *Uspacy) CreateCRMField(ctx context.Context, entityType string, fieldData any, opts ...RequestOption) (crm.Field, error) {
	responseBody, _, err := us.doPost(ctx, us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.CreateFieldUrl, entityType)), fieldData, opts...)
	return decodeJSON[crm.Field](responseBody, err)
}

// GetListValues returns the values of a CRM list field.
func (us *Uspacy) GetListValues(ctx context.Context, entityType, listName string, opts ...RequestOption) ([]crm.List, error) {
	responseBody, err := us.doGet(ctx, us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.ListsUrl, entityType, listName)), opts...)
	return decodeJSON[[]crm.List](responseBody, err)
}

// CreateListValues adds values to a CRM list field and returns the field's values.
func (us *Uspacy) CreateListValues(ctx context.Context, entityType, listName string, listValue any, opts ...RequestOption) ([]crm.List, error) {
	responseBody, _, err := us.doPost(ctx, us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.ListsUrl, entityType, listName)), listValue, opts...)
	return decodeJSON[[]crm.List](responseBody, err)
}

// CreateFailReasons returns all reasons for funnel with failWrite.ID
func (us *Uspacy) CreateFailReasons(ctx context.Context, failReason crm.Reason, opts ...RequestOption) (crm.Reason, error) {
	responseBody, _, err := us.doPost(ctx, us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.ReasonsUrl, failReason.ID)), crm.FailWrite{
		Title: failReason.Title,
		Sort:  failReason.Sort,
		Type:  "FAIL",
	}, opts...)
	return decodeJSON[crm.Reason](responseBody, err)
}

// CreateCall returns created call
func (us *Uspacy) CreateCall(ctx context.Context, callValue crm.Call, opts ...RequestOption) (crm.Call, error) {
	responseBody, _, err := us.doPost(ctx, us.buildURL(crm.VersionUrl, crm.CallUrl), callValue, opts...)
	return decodeJSON[crm.Call](responseBody, err)
}
