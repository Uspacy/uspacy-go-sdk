package api

import (
	"encoding/json"
	"uspacy-go-sdk/group"
)

// CreateTransferGroup creates a new transfer group
func (us *Uspacy) CreateTransferGroup(body interface{}) (group.TransferGroupOutput, error) {
	var groups group.TransferGroupOutput

	resp, err := us.doPostEmptyHeaders(buildURL(mainHost, group.VersionGroupUrl, group.TransferGroupsUrl), body)
	if err != nil {
		return groups, err
	}
	return groups, json.Unmarshal(resp, &groups)
}
