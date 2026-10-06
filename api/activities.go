package api

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/Uspacy/uspacy-go-sdk/v2/activities"
)

// CreateActivity sends a POST request to create a new activity using the provided entity data.
// It returns the created activity, the HTTP status code of the request, and any error encountered.
// If the request fails, the activity is the zero value.
func (us *Uspacy) CreateActivity(ctx context.Context, entityData map[string]any, opts ...RequestOption) (activities.Activity, int, error) {
	respBytes, code, err := us.doPost(ctx, us.buildURL(activities.VersionUrl, activities.ActivitiesUrl), entityData, opts...)
	v, err := decodeJSON[activities.Activity](respBytes, err)
	return v, code, err
}

// GetActivitiesList retrieves a list of activities based on the provided query parameters.
// It constructs the request URL using the base activities URL and optional query parameters (if provided).
// Returns an `ActivitiesList` containing the list of activities and any error encountered during the request or unmarshalling.
func (us *Uspacy) GetActivitiesList(ctx context.Context, params url.Values, opts ...RequestOption) (activities.ActivitiesList, error) {
	url := us.buildURL(activities.VersionUrl, activities.ActivitiesUrl)
	if len(params) != 0 {
		url = url + "?" + params.Encode()
	}
	body, err := us.doGet(ctx, url, opts...)
	return decodeJSON[activities.ActivitiesList](body, err)
}

// GetActivity retrieves details of a specific activity based on its entity ID and optional query parameters.
// The URL is constructed by formatting the activity URL with the given entity ID and appending query parameters if provided.
// Returns the requested `Activity` object and any error encountered during the request or unmarshalling.
func (us *Uspacy) GetActivity(ctx context.Context, entityId int64, params url.Values, opts ...RequestOption) (activities.Activity, error) {
	url := us.buildURL(activities.VersionUrl, fmt.Sprintf(activities.ActivityUrl, strconv.FormatInt(entityId, 10)))
	if len(params) != 0 {
		url = url + "?" + params.Encode()
	}
	body, err := us.doGet(ctx, url, opts...)
	return decodeJSON[activities.Activity](body, err)
}

// PatchActivity updates an existing activity identified by the entity ID with the provided entity data.
// It sends a PATCH request to the constructed URL.
// The function does not return any object, only an error if the request fails or if an issue occurs during the operation.
func (us *Uspacy) PatchActivity(ctx context.Context, entityId int64, entityData map[string]any, opts ...RequestOption) error {
	url := us.buildURL(activities.VersionUrl, fmt.Sprintf(activities.ActivityUrl, strconv.FormatInt(entityId, 10)))
	_, err := us.doPatch(ctx, url, entityData, opts...)
	if err != nil {
		return err
	}
	return nil
}

// DeleteActivity deletes an existing activity identified by the entity ID.
// It sends a DELETE request to the constructed URL and returns the HTTP status code and any error encountered during the request.
// The HTTP status code can be used to check if the deletion was successful.
func (us *Uspacy) DeleteActivity(ctx context.Context, entityId int64, opts ...RequestOption) (int, error) {
	url := us.buildURL(activities.VersionUrl, fmt.Sprintf(activities.ActivityUrl, strconv.FormatInt(entityId, 10)))
	return us.doDelete(ctx, url, nil, opts...)
}

// MassDeleteActivities sends a DELETE request to delete multiple activities based on the provided deletion data.
// The deletionData parameter contains the necessary information for mass deletion (like IDs or criteria).
// It constructs the request URL using the base mass deletion URL and sends the request with the provided body.
// Returns the HTTP status code of the request and any error encountered during the process.
//
// Renamed in v2 from MassDeletionActivies.
func (us *Uspacy) MassDeleteActivities(ctx context.Context, deletionData activities.MassDeletionBody, opts ...RequestOption) (int, error) {
	url := us.buildURL(activities.VersionUrl, activities.MassDeletion)
	return us.doDelete(ctx, url, deletionData, opts...)
}
