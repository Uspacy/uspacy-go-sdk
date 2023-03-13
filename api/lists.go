package api

import (
	"encoding/json"
	"fmt"

	"github.com/Uspacy/uspacy-go-sdk/crm"
)

// CreateListValues returns arrey of values for given type of CRM list
func (us *Uspacy) CreateListValues(entityType, listValue crm.Entity, fields interface{}) (createdLists []crm.List, err error) {
	responseBody, err := us.doPostEmptyHeaders(buildURL(mainHost, crm.VersionCRMUrl,
		fmt.Sprintf(crm.ListsUrl, entityType.GetUrl(crm.EntityType), listValue.GetUrl(crm.ListValue))), fields)
	if err != nil {
		return createdLists, err
	}
	return createdLists, json.Unmarshal(responseBody, &createdLists)
}
