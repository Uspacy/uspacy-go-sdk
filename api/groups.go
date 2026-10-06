package api

import (
	"context"
	"net/url"

	"github.com/Uspacy/uspacy-go-sdk/v2/group"
)

// GetGroups returns the groups matching params; pass nil for no filter.
func (us *Uspacy) GetGroups(ctx context.Context, params url.Values, opts ...RequestOption) (group.Groups, error) {
	urlStr := us.buildURL(group.VersionUrl, group.GroupUrl)
	if len(params) != 0 {
		urlStr += "?" + params.Encode()
	}
	body, err := us.doGet(ctx, urlStr, opts...)
	return decodeJSON[group.Groups](body, err)
}

// CreateGroup creates a group and returns it.
func (us *Uspacy) CreateGroup(ctx context.Context, groupData url.Values, opts ...RequestOption) (group.Group, error) {
	body, err := us.doPostEncodedForm(ctx, us.buildURL(group.VersionUrl, group.GroupUrl), groupData, opts...)
	return decodeJSON[group.Group](body, err)
}

// CreateTransferGroup creates groups in bulk through the transfer endpoint and returns the
// result.
func (us *Uspacy) CreateTransferGroup(ctx context.Context, body any, opts ...RequestOption) (group.TransferGroupOutput, error) {
	resp, _, err := us.doPost(ctx, us.buildURL(group.VersionUrl, group.TransferUrl), body, opts...)
	return decodeJSON[group.TransferGroupOutput](resp, err)
}
