package api

import (
	"encoding/json"
	"fmt"

	"github.com/Uspacy/uspacy-go-sdk/crm"
)

// CreateFunnel returns created funnel
func (us *Uspacy) CreateFunnel(entity string, body interface{}) (object crm.Funnel, err error) {
	responseBody, err := us.doPostEmptyHeaders(buildURL(mainHost, crm.VersionUrl, fmt.Sprintf(crm.CreateFunnelUrl, entity)), body)
	if err != nil {
		return object, err
	}
	return object, json.Unmarshal(responseBody, &object)
}

// CreateKanbanStage returns created kanban stage
func (us *Uspacy) CreateKanbanStage(entity string, body interface{}) (object crm.KanbanStage, err error) {
	responseBody, err := us.doPostEmptyHeaders(buildURL(mainHost, crm.VersionUrl, fmt.Sprintf(crm.CreateFunnelUrl, entity)), body)
	if err != nil {
		return object, err
	}
	return object, json.Unmarshal(responseBody, &object)
}

// CreateContact returns created contacts object
func (us *Uspacy) CreateContact(entityType crm.Entity, entity map[string]interface{}) (object crm.Contacts, err error) {
	body, err := us.doPostEmptyHeaders(buildURL(mainHost, crm.VersionUrl, fmt.Sprintf(crm.CreateEntity, entityType.GetUrl())), entity)
	if err != nil {
		return object, err
	}
	return object, json.Unmarshal(body, &object)
}

// CreateCompany returns created contacts object
func (us *Uspacy) CreateCompany(entityType crm.Entity, entity map[string]interface{}) (object crm.Companies, err error) {
	body, err := us.doPostEmptyHeaders(buildURL(mainHost, crm.VersionUrl, fmt.Sprintf(crm.CreateEntity, entityType.GetUrl())), entity)
	if err != nil {
		return object, err
	}
	return object, json.Unmarshal(body, &object)
}

// CreateLeads returns created contacts object
func (us *Uspacy) CreateLeads(entityType crm.Entity, entity map[string]interface{}) (object crm.Leads, err error) {
	body, err := us.doPostEmptyHeaders(buildURL(mainHost, crm.VersionUrl, fmt.Sprintf(crm.CreateEntity, entityType.GetUrl())), entity)
	if err != nil {
		return object, err
	}
	return object, json.Unmarshal(body, &object)
}

// CreateDeals returns created contacts object
func (us *Uspacy) CreateDeals(entityType crm.Entity, entity map[string]interface{}) (object crm.Deals, err error) {
	body, err := us.doPostEmptyHeaders(buildURL(mainHost, crm.VersionUrl, fmt.Sprintf(crm.CreateEntity, entityType.GetUrl())), entity)
	if err != nil {
		return object, err
	}
	return object, json.Unmarshal(body, &object)
}
