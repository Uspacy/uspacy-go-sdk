package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/Uspacy/uspacy-go-sdk/v2/crm"
)

// CreateEntity this method does not return any object, just error
func (us *Uspacy) CreateEntity(ctx context.Context, entityType string, entityData map[string]any, opts ...RequestOption) (int64, int, error) {
	respBytes, code, err := us.doPost(ctx, us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.EntityUrl, entityType)), entityData, opts...)
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

// GetCrmEntitiesList this method return arrey of entities present in crm and error
func (us *Uspacy) GetCrmEntitiesList(ctx context.Context) (entities []crm.CrmEntities, err error) {
	body, err := us.doGetEmptyHeaders(ctx, us.buildURL(crm.VersionUrl, crm.EntitiesUrl))
	if err != nil {
		return entities, err
	}
	var resp = crm.CrmEntitiesList{}
	err = json.Unmarshal(body, &resp)
	return resp.Data, err
}

// GetEntities this method return arrey of entities present in crm and error
func (us *Uspacy) GetEntities(ctx context.Context, entityType string, params url.Values) (entities crm.CRMEntity, err error) {
	body, err := us.doGetEmptyHeaders(ctx, us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.EntityUrl, entityType))+"?"+params.Encode())
	if err != nil {
		return entities, err
	}
	return entities, json.Unmarshal(body, &entities)
}

// GetCRMEntitiesForExport this method return arrey of objects wifh all fields and error
func (us *Uspacy) GetCRMEntitiesForExport(ctx context.Context, entityType string, params url.Values) (entities crm.CRMEntityForExport, err error) {
	var entityRoute string
	switch entityType {
	case crm.ProductsNum.GetUrl():
		entityRoute = fmt.Sprintf(crm.ProductsUrl, "")
	default:
		entityRoute = fmt.Sprintf(crm.EntityUrl, entityType)
	}
	body, err := us.doGetEmptyHeaders(ctx, us.buildURL(crm.VersionUrl, entityRoute)+"?"+params.Encode())
	if err != nil {
		return entities, err
	}
	return entities, json.Unmarshal(body, &entities)
}

// GetContacts returns an array of contact objects and an error.
func (us *Uspacy) GetContacts(ctx context.Context, params url.Values) (entities crm.Contacts, err error) {
	body, err := us.doGetEmptyHeaders(ctx, us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.EntityUrl, crm.ContactsNum.GetUrl()))+"?"+params.Encode())
	if err != nil {
		return entities, err
	}
	return entities, json.Unmarshal(body, &entities)
}

// GetDeals returns an array of deal objects and an error.
func (us *Uspacy) GetDeals(ctx context.Context, params url.Values) (entities crm.Deals, err error) {
	body, err := us.doGetEmptyHeaders(ctx, us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.EntityUrl, crm.DealsNum.GetUrl()))+"?"+params.Encode())
	if err != nil {
		return entities, err
	}
	return entities, json.Unmarshal(body, &entities)
}

// GetLeads returns an array of lead objects and an error.
func (us *Uspacy) GetLeads(ctx context.Context, params url.Values) (entities crm.Leads, err error) {
	body, err := us.doGetEmptyHeaders(ctx, us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.EntityUrl, crm.LeadsNum.GetUrl()))+"?"+params.Encode())
	if err != nil {
		return entities, err
	}
	return entities, json.Unmarshal(body, &entities)
}

// GetList returns raw response for CRM entities with filters
// This method is useful for searching entities with custom filters
// entityType should be one of: crm.LeadsNum.GetUrl(), crm.DealsNum.GetUrl(), crm.ContactsNum.GetUrl(), crm.CompaniesNum.GetUrl()
func (us *Uspacy) GetList(ctx context.Context, entityType string, params url.Values, opts ...RequestOption) ([]byte, error) {
	url := us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.EntityUrl, entityType))
	if len(params) > 0 {
		url += "?" + params.Encode()
	}

	return us.doGetEmptyHeaders(ctx, url, opts...)
}

// GetEntity returns one record as raw JSON: GET crm/v1/entities/{entityType}/{id}.
// The id is its own path segment so the leading slash is preserved.
// A non-2xx response is returned as *HTTPError.
func (us *Uspacy) GetEntity(ctx context.Context, entityType string, id int64) ([]byte, error) {
	url := us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.EntityUrl, entityType), strconv.FormatInt(id, 10))
	body, _, err := us.doRaw(ctx, url, http.MethodGet, headersMap, nil)
	return body, err
}

// PatchEntity this method does not return any object, just error
func (us *Uspacy) PatchEntity(ctx context.Context, entityType string, id string, entityData map[string]any) error {
	_, err := us.doPatchEmptyHeaders(ctx, us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.EntityUrl, entityType), id), entityData)
	if err != nil {
		return err
	}
	return nil
}

// EntityMassEdit this method does not return any object, just error
func (us *Uspacy) EntityMassEdit(ctx context.Context, entityType string, entityData map[string]any) error {
	_, err := us.doPatchEmptyHeaders(ctx, us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.EntityUrl, entityType), "mass_edit"), entityData)
	if err != nil {
		return err
	}
	return nil
}

// CreateContact returns created contact object
func (us *Uspacy) CreateContact(ctx context.Context, contactData map[string]any, opts ...RequestOption) (contact crm.Contact, err error) {
	body, _, err := us.doPost(ctx, us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.EntityUrl, crm.ContactsNum.GetUrl())), contactData, opts...)
	if err != nil {
		return contact, err
	}
	return contact, json.Unmarshal(body, &contact)
}

// CreateCompany returns created company object
func (us *Uspacy) CreateCompany(ctx context.Context, companyData map[string]any, opts ...RequestOption) (company crm.Company, err error) {
	body, _, err := us.doPost(ctx, us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.EntityUrl, crm.CompaniesNum.GetUrl())), companyData, opts...)
	if err != nil {
		return company, err
	}
	return company, json.Unmarshal(body, &company)
}

// CreateLeads returns created lead object
func (us *Uspacy) CreateLead(ctx context.Context, leadData map[string]any, opts ...RequestOption) (lead crm.Lead, err error) {
	body, _, err := us.doPost(ctx, us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.EntityUrl, crm.LeadsNum.GetUrl())), leadData, opts...)
	if err != nil {
		return lead, err
	}
	return lead, json.Unmarshal(body, &lead)
}

// CreateDeals returns created deal object
func (us *Uspacy) CreateDeal(ctx context.Context, dealData map[string]any, opts ...RequestOption) (deal crm.Deal, err error) {
	body, _, err := us.doPost(ctx, us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.EntityUrl, crm.DealsNum.GetUrl())), dealData, opts...)
	if err != nil {
		return deal, err
	}
	return deal, json.Unmarshal(body, &deal)
}

// GetField returns Field struct for a given type of entity & field
func (us *Uspacy) GetField(ctx context.Context, entityType string, fieldType string) (field crm.Field, err error) {
	body, err := us.doGetEmptyHeaders(ctx, us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.FieldsUrl, entityType, fieldType)))
	if err != nil {
		return field, err
	}
	return field, json.Unmarshal(body, &field)
}

// DeleteField delete selected field for given type of entity
func (us *Uspacy) DeleteField(ctx context.Context, entityType string, codeField string) (err error) {
	_, err = us.doDeleteEmptyHeaders(ctx, us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.FieldsUrl, entityType, codeField)), nil)
	return err
}

// GetFields returns the fields of an entity type: GET crm/v1/entities/{entityType}/fields.
// A non-2xx response is returned as *HTTPError.
func (us *Uspacy) GetFields(ctx context.Context, entityType string) ([]crm.Field, error) {
	url := us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.FieldsUrl, entityType, ""))
	body, _, err := us.doRaw(ctx, url, http.MethodGet, headersMap, nil)
	if err != nil {
		return nil, err
	}
	var resp crm.Fields
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// CreateFunnel returns created funnel
func (us *Uspacy) CreateFunnel(ctx context.Context, entityType string, funnelData any, opts ...RequestOption) (entityFunnel crm.Funnel, err error) {
	responseBody, _, err := us.doPost(ctx, us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.FunnelUrl, entityType)), funnelData, opts...)
	if err != nil {
		return entityFunnel, err
	}
	return entityFunnel, json.Unmarshal(responseBody, &entityFunnel)
}

// GetFunnels returns funnels by entityType
func (us *Uspacy) GetFunnels(ctx context.Context, entityType string) (funnels crm.FunnelsById, err error) {
	responseBody, err := us.doGetEmptyHeaders(ctx, us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.FunnelUrl, entityType)))
	if err != nil {
		return funnels, err
	}
	return funnels, json.Unmarshal(responseBody, &funnels)
}

// CreateFunnelStage returns created kanban stage
func (us *Uspacy) CreateFunnelStage(ctx context.Context, entityType string, stageData any, opts ...RequestOption) (kanbanStage crm.KanbanStage, err error) {
	responseBody, _, err := us.doPost(ctx, us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.KanbanStageUrl, entityType, "")), stageData, opts...)
	if err != nil {
		return kanbanStage, err
	}
	return kanbanStage, json.Unmarshal(responseBody, &kanbanStage)
}

// GetAllFunnelStages returns all kanban stages
func (us *Uspacy) GetAllFunnelStages(ctx context.Context, entityType string) (kanbanStages []crm.KanbanStage, err error) {
	responseBody, err := us.doGetEmptyHeaders(ctx, us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.KanbanStageUrl, entityType, "")))
	if err != nil {
		return kanbanStages, err
	}
	var resp crm.KanbanStages
	if err := json.Unmarshal(responseBody, &resp); err != nil {
		return kanbanStages, err
	}
	return resp.Data, nil
}

// GetFunnelStageDyId returns kanban stage
func (us *Uspacy) GetFunnelStageDyId(ctx context.Context, entityType string, id int) (kanbanStages crm.KanbanStages, err error) {
	responseBody, err := us.doGetEmptyHeaders(ctx, us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.KanbanStageUrl, entityType, fmt.Sprintf(crm.StageByFunnelIdUrl, id))))
	if err != nil {
		return kanbanStages, err
	}
	return kanbanStages, json.Unmarshal(responseBody, &kanbanStages)
}

// PatchFunnelStage returns kanban stage
func (us *Uspacy) PatchFunnelStage(ctx context.Context, entityType string, id int, stage crm.FunnelStage) (kanbanStage crm.KanbanStage, err error) {
	responseBody, err := us.doPatchEmptyHeaders(ctx, us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.KanbanStageUrl, entityType, id)), stage)
	if err != nil {
		return kanbanStage, err
	}
	return kanbanStage, json.Unmarshal(responseBody, &kanbanStage)
}

// MoveFunnelStage moves a funnel stage
func (us *Uspacy) MoveFunnelStage(ctx context.Context, entityType string, entityId int64, stageId string, reason crm.KanbanFailReasonCRM, opts ...RequestOption) (err error) {
	_, _, err = us.doPost(ctx, us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.MoveKanbanStageUrl, entityType, entityId, stageId)), reason, opts...)
	return err
}

// CreateCRMField in CRM entity returns created field
func (us *Uspacy) CreateCRMField(ctx context.Context, entityType string, fieldData any, opts ...RequestOption) (entityField crm.Field, err error) {
	responseBody, _, err := us.doPost(ctx, us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.CreateFieldUrl, entityType)), fieldData, opts...)
	if err != nil {
		return entityField, err
	}
	return entityField, json.Unmarshal(responseBody, &entityField)
}

// GetListValues returns arrey of values for given type of CRM list
func (us *Uspacy) GetListValues(ctx context.Context, entityType, listName string) (lists []crm.List, err error) {
	responseBody, err := us.doGetEmptyHeaders(ctx, us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.ListsUrl, entityType, listName)))
	if err != nil {
		return lists, err
	}
	return lists, json.Unmarshal(responseBody, &lists)
}

// CreateListValues returns arrey of values for given type of CRM list
func (us *Uspacy) CreateListValues(ctx context.Context, entityType, listName string, listValue any, opts ...RequestOption) (lists []crm.List, err error) {
	responseBody, _, err := us.doPost(ctx, us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.ListsUrl, entityType, listName)), listValue, opts...)
	if err != nil {
		return lists, err
	}
	return lists, json.Unmarshal(responseBody, &lists)
}

// CreateFailReasons returns all reasons for funnel with failWrite.ID
func (us *Uspacy) CreateFailReasons(ctx context.Context, failReason crm.Reason, opts ...RequestOption) (reasons crm.Reason, err error) {
	responseBody, _, err := us.doPost(ctx, us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.ReasonsUrl, failReason.ID)), crm.FailWrite{
		Title: failReason.Title,
		Sort:  failReason.Sort,
		Type:  "FAIL",
	}, opts...)
	if err != nil {
		return reasons, err
	}
	return reasons, json.Unmarshal(responseBody, &reasons)
}

// CreateCall returns created call
func (us *Uspacy) CreateCall(ctx context.Context, callValue crm.Call, opts ...RequestOption) (call crm.Call, err error) {
	responseBody, _, err := us.doPost(ctx, us.buildURL(crm.VersionUrl, crm.CallUrl), callValue, opts...)
	if err != nil {
		return call, err
	}
	return call, json.Unmarshal(responseBody, &call)
}
