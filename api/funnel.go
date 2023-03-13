package api

import (
	"encoding/json"
	"fmt"

	"github.com/Uspacy/uspacy-go-sdk/crm"
)

// CreateFunnel returns created funnel
func (us *Uspacy) CreateFunnel(entity string, body interface{}) (funnel crm.Funnel, err error) {
	responseBody, err := us.doPostEmptyHeaders(buildURL(mainHost, crm.VersionCRMUrl, fmt.Sprintf(crm.FunnelUrl, entity)), body)
	if err != nil {
		return funnel, err
	}
	return funnel, json.Unmarshal(responseBody, &funnel)
}
