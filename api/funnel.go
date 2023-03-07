package api

import (
	"encoding/json"
	"fmt"
	"io"
	"uspacy-go-sdk/crm"
)

// CreateFunnel created funnel
func (us *Uspacy) CreateFunnel(entity string, body io.Reader) (crm.Funnel, error) {
	var funnel crm.Funnel

	requestBody, err := us.doPostEmptyHeaders(buildURL(mainHost, crm.VersionUrl, fmt.Sprintf(crm.CreateFunnelUrl, entity)), body)
	if err != nil {
		return funnel, err
	}
	return funnel, json.Unmarshal(requestBody, &funnel)
}
