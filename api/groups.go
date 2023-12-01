package api

import (
	"encoding/json"
	"net/url"

	"github.com/Uspacy/uspacy-go-sdk/group"
)

// GetGroups returns  list of groups
func (us *Uspacy) GetGroups() (groups group.Groups, err error) {
	body, err := us.doGetEmptyHeaders(us.buildURL(group.VersionUrl, group.GroupUrl))
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
func (us *Uspacy) CreateTransferGroup(body interface{}) (groups group.TransferGroupOutput, err error) {
	resp, _, err := us.doPostEmptyHeaders(us.buildURL(group.VersionUrl, group.TransferUrl), body)
	if err != nil {
		return groups, err
	}
	return groups, json.Unmarshal(resp, &groups)
}
