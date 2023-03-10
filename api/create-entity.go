package api

import (
	"encoding/json"
	"fmt"

	"github.com/Uspacy/uspacy-go-sdk/crm"
)

func (us *Uspacy) CreateEntity(entityType crm.Entity, entity map[string]interface{}) (object crm.EntityResponse, err error) {
	body, err := us.doPostEmptyHeaders(buildURL(mainHost, crm.VersionUrl, fmt.Sprintf(crm.EntityUrl, entityType.GetUrl(crm.EntityType))), entity)
	if err != nil {
		return object, err
	}
	return object, json.Unmarshal(body, &object)
}
