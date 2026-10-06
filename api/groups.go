package api

import (
	"context"
	"net/url"

	"github.com/Uspacy/uspacy-go-sdk/v2/group"
)

// GetGroups returns  list of groups
func (us *Uspacy) GetGroups(ctx context.Context, params ...url.Values) (group.Groups, error) {
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
	body, err := us.doGetEmptyHeaders(ctx, urlStr)
	return decodeJSON[group.Groups](body, err)
}

// CreateGroup returns created group object
func (us *Uspacy) CreateGroup(ctx context.Context, groupData url.Values) (group.Group, error) {
	body, err := us.doPostEncodedForm(ctx, us.buildURL(group.VersionUrl, group.GroupUrl), groupData)
	return decodeJSON[group.Group](body, err)
}

// CreateTransferGroup creates a new transfer group
func (us *Uspacy) CreateTransferGroup(ctx context.Context, body any, opts ...RequestOption) (group.TransferGroupOutput, error) {
	resp, _, err := us.doPost(ctx, us.buildURL(group.VersionUrl, group.TransferUrl), body, opts...)
	return decodeJSON[group.TransferGroupOutput](resp, err)
}
