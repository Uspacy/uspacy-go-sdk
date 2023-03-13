package api

import (
	"encoding/json"
	"fmt"
	"github.com/Uspacy/uspacy-go-sdk/crm"
)

// CreateFunnel returns created funnel
func (us *Uspacy) CreateFunnel(entity string, body interface{}) (funnel crm.Funnel, err error) {
	responseBody, err := us.doPostEmptyHeaders(buildURL(mainHost, crm.VersionUrl, fmt.Sprintf(crm.CreateFunnelUrl, entity)), body)
	if err != nil {
		return funnel, err
	}
	return funnel, json.Unmarshal(responseBody, &funnel)
}

// CreateKanbanStage returns created kanban stage
func (us *Uspacy) CreateKanbanStage(entity string, body interface{}) (kanbanStage crm.KanbanStage, err error) {
	responseBody, err := us.doPostEmptyHeaders(buildURL(mainHost, crm.VersionUrl, fmt.Sprintf(crm.CreateFunnelUrl, entity)), body)
	if err != nil {
		return kanbanStage, err
	}
	return kanbanStage, json.Unmarshal(responseBody, &kanbanStage)
}

// CreateContact returns created contacts object
func (us *Uspacy) CreateContact(entityType crm.Entity, entity map[string]interface{}) (object crm.ContactsRes, err error) {
	body, err := us.doPostEmptyHeaders(buildURL(mainHost, crm.VersionUrl, fmt.Sprintf(crm.CreateEntity, entityType.GetUrl())), entity)
	if err != nil {
		return object, err
	}
	return object, json.Unmarshal(body, &object)
}

// CreateCompany returns created contacts object
func (us *Uspacy) CreateCompany(entityType crm.Entity, entity map[string]interface{}) (object crm.CompaniesRes, err error) {
	body, err := us.doPostEmptyHeaders(buildURL(mainHost, crm.VersionUrl, fmt.Sprintf(crm.CreateEntity, entityType.GetUrl())), entity)
	if err != nil {
		return object, err
	}
	return object, json.Unmarshal(body, &object)
}

// CreateLeads returns created contacts object
func (us *Uspacy) CreateLeads(entityType crm.Entity, entity map[string]interface{}) (object crm.LeadsRes, err error) {
	body, err := us.doPostEmptyHeaders(buildURL(mainHost, crm.VersionUrl, fmt.Sprintf(crm.CreateEntity, entityType.GetUrl())), entity)
	if err != nil {
		return object, err
	}
	return object, json.Unmarshal(body, &object)
}

// CreateDeals returns created contacts object
func (us *Uspacy) CreateDeals(entityType crm.Entity, entity map[string]interface{}) (object crm.DealsRes, err error) {
	body, err := us.doPostEmptyHeaders(buildURL(mainHost, crm.VersionUrl, fmt.Sprintf(crm.CreateEntity, entityType.GetUrl())), entity)
	if err != nil {
		return object, err
	}
	return object, json.Unmarshal(body, &object)
}
