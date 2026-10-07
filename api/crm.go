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
	body, err := us.doGet(ctx, withQuery(us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.EntityUrl, entityType)), params), opts...)
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
	body, err := us.doGet(ctx, withQuery(us.buildURL(crm.VersionUrl, entityRoute), params), opts...)
	return decodeJSON[crm.CRMEntityForExport](body, err)
}

// GetContacts returns the contacts matching params.
func (us *Uspacy) GetContacts(ctx context.Context, params url.Values, opts ...RequestOption) (crm.Contacts, error) {
	body, err := us.doGet(ctx, withQuery(us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.EntityUrl, crm.ContactsNum.GetUrl())), params), opts...)
	return decodeJSON[crm.Contacts](body, err)
}

// GetDeals returns the deals matching params.
func (us *Uspacy) GetDeals(ctx context.Context, params url.Values, opts ...RequestOption) (crm.Deals, error) {
	body, err := us.doGet(ctx, withQuery(us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.EntityUrl, crm.DealsNum.GetUrl())), params), opts...)
	return decodeJSON[crm.Deals](body, err)
}

// GetLeads returns the leads matching params.
func (us *Uspacy) GetLeads(ctx context.Context, params url.Values, opts ...RequestOption) (crm.Leads, error) {
	body, err := us.doGet(ctx, withQuery(us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.EntityUrl, crm.LeadsNum.GetUrl())), params), opts...)
	return decodeJSON[crm.Leads](body, err)
}

// GetList returns the records of a CRM entity type matching params as raw JSON, for
// filters the typed methods do not cover. entityType is one of crm.LeadsNum.GetUrl(),
// crm.DealsNum.GetUrl(), crm.ContactsNum.GetUrl() or crm.CompaniesNum.GetUrl().
func (us *Uspacy) GetList(ctx context.Context, entityType string, params url.Values, opts ...RequestOption) ([]byte, error) {
	url := withQuery(us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.EntityUrl, entityType)), params)

	return us.doGet(ctx, url, opts...)
}

// GetEntity returns a CRM record as raw JSON.
func (us *Uspacy) GetEntity(ctx context.Context, entityType string, id int64, opts ...RequestOption) ([]byte, error) {
	return us.doGet(ctx, us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.EntityUrl, entityType), strconv.FormatInt(id, 10)), opts...)
}

// PatchEntity updates a CRM record.
func (us *Uspacy) PatchEntity(ctx context.Context, entityType string, id string, entityData map[string]any, opts ...RequestOption) error {
	_, err := us.doPatch(ctx, us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.EntityUrl, entityType), id), entityData, opts...)
	return err
}

// EntityMassEdit updates several CRM records of one entity type at once.
func (us *Uspacy) EntityMassEdit(ctx context.Context, entityType string, entityData map[string]any, opts ...RequestOption) error {
	_, err := us.doPatch(ctx, us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.EntityUrl, entityType), "mass_edit"), entityData, opts...)
	return err
}

// CreateContact creates a contact and returns it.
func (us *Uspacy) CreateContact(ctx context.Context, contactData map[string]any, opts ...RequestOption) (crm.Contact, error) {
	body, _, err := us.doPost(ctx, us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.EntityUrl, crm.ContactsNum.GetUrl())), contactData, opts...)
	return decodeJSON[crm.Contact](body, err)
}

// CreateCompany creates a company and returns it.
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

// GetField returns the field with code fieldType of a CRM entity type.
func (us *Uspacy) GetField(ctx context.Context, entityType string, fieldType string, opts ...RequestOption) (crm.Field, error) {
	body, err := us.doGet(ctx, us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.FieldsUrl, entityType, fieldType)), opts...)
	return decodeJSON[crm.Field](body, err)
}

// DeleteField deletes a field of a CRM entity type.
func (us *Uspacy) DeleteField(ctx context.Context, entityType string, codeField string, opts ...RequestOption) (err error) {
	_, err = us.doDelete(ctx, us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.FieldsUrl, entityType, codeField)), nil, opts...)
	return err
}

// GetFields returns the fields of a CRM entity type.
func (us *Uspacy) GetFields(ctx context.Context, entityType string, opts ...RequestOption) ([]crm.Field, error) {
	body, err := us.doGet(ctx, us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.FieldsUrl, entityType, "")), opts...)
	resp, err := decodeJSON[crm.Fields](body, err)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// CreateFunnel creates a funnel for a CRM entity type and returns it.
func (us *Uspacy) CreateFunnel(ctx context.Context, entityType string, funnelData any, opts ...RequestOption) (crm.Funnel, error) {
	responseBody, _, err := us.doPost(ctx, us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.FunnelUrl, entityType)), funnelData, opts...)
	return decodeJSON[crm.Funnel](responseBody, err)
}

// GetFunnels returns the funnels of a CRM entity type.
func (us *Uspacy) GetFunnels(ctx context.Context, entityType string, opts ...RequestOption) (crm.FunnelsById, error) {
	responseBody, err := us.doGet(ctx, us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.FunnelUrl, entityType)), opts...)
	return decodeJSON[crm.FunnelsById](responseBody, err)
}

// CreateFunnelStage creates a kanban stage for a CRM entity type and returns it.
func (us *Uspacy) CreateFunnelStage(ctx context.Context, entityType string, stageData any, opts ...RequestOption) (crm.KanbanStage, error) {
	responseBody, _, err := us.doPost(ctx, us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.KanbanStageUrl, entityType, "")), stageData, opts...)
	return decodeJSON[crm.KanbanStage](responseBody, err)
}

// GetAllFunnelStages returns all kanban stages of a CRM entity type.
func (us *Uspacy) GetAllFunnelStages(ctx context.Context, entityType string, opts ...RequestOption) ([]crm.KanbanStage, error) {
	responseBody, err := us.doGet(ctx, us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.KanbanStageUrl, entityType, "")), opts...)
	resp, err := decodeJSON[crm.KanbanStages](responseBody, err)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// GetFunnelStageById returns the kanban stages of the funnel with the given id.
//
// Renamed in v2 from GetFunnelStageDyId.
func (us *Uspacy) GetFunnelStageById(ctx context.Context, entityType string, id int, opts ...RequestOption) (crm.KanbanStages, error) {
	responseBody, err := us.doGet(ctx, us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.KanbanStageUrl, entityType, fmt.Sprintf(crm.StageByFunnelIdUrl, id))), opts...)
	return decodeJSON[crm.KanbanStages](responseBody, err)
}

// PatchFunnelStage updates a kanban stage and returns it.
func (us *Uspacy) PatchFunnelStage(ctx context.Context, entityType string, id int, stage crm.FunnelStage, opts ...RequestOption) (crm.KanbanStage, error) {
	responseBody, err := us.doPatch(ctx, us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.KanbanStageUrl, entityType, id)), stage, opts...)
	return decodeJSON[crm.KanbanStage](responseBody, err)
}

// MoveFunnelStage moves a CRM record to another kanban stage. reason explains a move to
// a failure stage.
func (us *Uspacy) MoveFunnelStage(ctx context.Context, entityType string, entityId int64, stageId string, reason crm.KanbanFailReasonCRM, opts ...RequestOption) (err error) {
	_, _, err = us.doPost(ctx, us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.MoveKanbanStageUrl, entityType, entityId, stageId)), reason, opts...)
	return err
}

// CreateCRMField creates a field of a CRM entity type and returns it.
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

// CreateFailReasons adds a failure reason to the funnel failReason.ID and returns it.
func (us *Uspacy) CreateFailReasons(ctx context.Context, failReason crm.Reason, opts ...RequestOption) (crm.Reason, error) {
	responseBody, _, err := us.doPost(ctx, us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.ReasonsUrl, failReason.ID)), crm.FailWrite{
		Title: failReason.Title,
		Sort:  failReason.Sort,
		Type:  "FAIL",
	}, opts...)
	return decodeJSON[crm.Reason](responseBody, err)
}

// CreateCall records a call and returns it.
func (us *Uspacy) CreateCall(ctx context.Context, callValue crm.Call, opts ...RequestOption) (crm.Call, error) {
	responseBody, _, err := us.doPost(ctx, us.buildURL(crm.VersionUrl, crm.CallUrl), callValue, opts...)
	return decodeJSON[crm.Call](responseBody, err)
}
