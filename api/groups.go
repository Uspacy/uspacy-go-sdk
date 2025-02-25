package api

import (
	"encoding/json"
	"net/url"

	"github.com/Uspacy/uspacy-go-sdk/group"
)

// GetGroups returns  list of groups
func (us *Uspacy) GetGroups(params ...url.Values) (groups group.Groups, err error) {
	urlStr := us.buildURL(group.VersionUrl, group.GroupUrl)
	if len(params) != 0 {
		mergedParams := make(url.Values)
		for _, p := range params {
			for key, values := range p {
				for _, value := range values {
					mergedParams.Add(key, value)
				}
			}
		}
		urlStr = urlStr + "?" + mergedParams.Encode()
	}
	body, err := us.doGetEmptyHeaders(urlStr)
	if err != nil {
		return groups, err
	}
	return groups, json.Unmarshal(body, &groups)
}

// CreateGroup returns created group object
func (us *Uspacy) CreateGroup(groupData url.Values) (_group group.Group, err error) {
	body, err := us.doPostEncodedForm(us.buildURL(group.VersionUrl, group.GroupUrl), groupData)
	if err != nil {
		return _group, err
	}
	return _group, json.Unmarshal(body, &_group)
}

// CreateTransferGroup creates a new transfer group
func (us *Uspacy) CreateTransferGroup(body interface{}, headers ...map[string]string) (groups group.TransferGroupOutput, err error) {
	resp, _, err := us.doPost(us.buildURL(group.VersionUrl, group.TransferUrl), body, headers...)
	if err != nil {
		return groups, err
	}
	return groups, json.Unmarshal(resp, &groups)
}
