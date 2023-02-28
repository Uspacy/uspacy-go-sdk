package api

import (
	"encoding/json"
	"uspacy-go-sdk/group"
)

// CreateTransferGroup creates a new transfer group
func (us *Uspacy) CreateTransferGroup() (group.TransferGroupOutput, error) {
	var groups group.TransferGroupOutput

	body, err := us.doGetEmptyHeaders(buildURL(mainHost, group.VersionGroupUrl, group.GroupsUrl))
	if err != nil {
		return groups, err
	}
	return groups, json.Unmarshal(body, &groups)
}
