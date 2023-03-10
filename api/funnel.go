package api

import (
	"encoding/json"
	"fmt"

	"github.com/Uspacy/uspacy-go-sdk/crm"
)

// CreateFunnel returns created funnel
func (us *Uspacy) CreateFunnel(entity string, body interface{}) (crm.Funnel, error) {
	var funnel crm.Funnel

	responseBody, err := us.doPostEmptyHeaders(buildURL(mainHost, crm.VersionUrl, fmt.Sprintf(crm.CreateFunnelUrl, entity)), body)
	if err != nil {
		return funnel, err
	}
	return funnel, json.Unmarshal(responseBody, &funnel)
}
