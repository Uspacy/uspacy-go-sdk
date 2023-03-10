package api

import (
	"encoding/json"

	"github.com/Uspacy/uspacy-go-sdk/group"
)

// CreateTransferGroup creates a new transfer group
func (us *Uspacy) CreateTransferGroup(body interface{}) (groups group.TransferGroupOutput, err error) {
	resp, err := us.doPostEmptyHeaders(buildURL(mainHost, group.VersionGroupUrl, group.TransferGroupsUrl), body)
	if err != nil {
		return groups, err
	}
	return groups, json.Unmarshal(resp, &groups)
}
