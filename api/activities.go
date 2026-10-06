package api

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/Uspacy/uspacy-go-sdk/v2/activities"
)

// CreateActivity creates an activity and returns it with the HTTP status code.
func (us *Uspacy) CreateActivity(ctx context.Context, entityData map[string]any, opts ...RequestOption) (activities.Activity, int, error) {
	respBytes, code, err := us.doPost(ctx, us.buildURL(activities.VersionUrl, activities.ActivitiesUrl), entityData, opts...)
	v, err := decodeJSON[activities.Activity](respBytes, err)
	return v, code, err
}

// GetActivitiesList returns the activities matching params; pass nil for no filter.
func (us *Uspacy) GetActivitiesList(ctx context.Context, params url.Values, opts ...RequestOption) (activities.ActivitiesList, error) {
	url := us.buildURL(activities.VersionUrl, activities.ActivitiesUrl)
	if len(params) != 0 {
		url = url + "?" + params.Encode()
	}
	body, err := us.doGet(ctx, url, opts...)
	return decodeJSON[activities.ActivitiesList](body, err)
}

// GetActivity returns an activity by ID. params is optional; pass nil for none.
func (us *Uspacy) GetActivity(ctx context.Context, entityId int64, params url.Values, opts ...RequestOption) (activities.Activity, error) {
	url := us.buildURL(activities.VersionUrl, fmt.Sprintf(activities.ActivityUrl, strconv.FormatInt(entityId, 10)))
	if len(params) != 0 {
		url = url + "?" + params.Encode()
	}
	body, err := us.doGet(ctx, url, opts...)
	return decodeJSON[activities.Activity](body, err)
}

// PatchActivity updates an activity.
func (us *Uspacy) PatchActivity(ctx context.Context, entityId int64, entityData map[string]any, opts ...RequestOption) error {
	url := us.buildURL(activities.VersionUrl, fmt.Sprintf(activities.ActivityUrl, strconv.FormatInt(entityId, 10)))
	_, err := us.doPatch(ctx, url, entityData, opts...)
	if err != nil {
		return err
	}
	return nil
}

// DeleteActivity deletes an activity and returns the HTTP status code.
func (us *Uspacy) DeleteActivity(ctx context.Context, entityId int64, opts ...RequestOption) (int, error) {
	url := us.buildURL(activities.VersionUrl, fmt.Sprintf(activities.ActivityUrl, strconv.FormatInt(entityId, 10)))
	return us.doDelete(ctx, url, nil, opts...)
}

// MassDeleteActivities deletes the activities selected by deletionData and returns the HTTP
// status code.
//
// Renamed in v2 from MassDeletionActivies.
func (us *Uspacy) MassDeleteActivities(ctx context.Context, deletionData activities.MassDeletionBody, opts ...RequestOption) (int, error) {
	url := us.buildURL(activities.VersionUrl, activities.MassDeletion)
	return us.doDelete(ctx, url, deletionData, opts...)
}
