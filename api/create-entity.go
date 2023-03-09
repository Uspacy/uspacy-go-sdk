package api

import (
	"encoding/json"
	"fmt"
	"github.com/Uspacy/uspacy-go-sdk/crm"
)

func (us *Uspacy) CreateEntity(entityType crm.Entity, entity map[string]interface{}) (object crm.EntityResponse, err error) {
	body, err := us.doPostEmptyHeaders(buildURL(mainHost, crm.VersionUrl, fmt.Sprintf(crm.CreateEntity, entityType.GetUrl())), entity)
	return object, json.Unmarshal(body, &object)
}
