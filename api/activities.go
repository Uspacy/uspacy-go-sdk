package api

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/Uspacy/uspacy-go-sdk/activities"
)

// CreateActivity sends a POST request to create a new activity using the provided entity data.
// It returns the created activity's ID, the HTTP status code of the request, and any error encountered.
// If an error occurs during the request or while unmarshalling the response, the error is returned along with a zero value for the ID.
func (us *Uspacy) CreateActivity(entityData map[string]interface{}) (entity activities.Activity, code int, err error) {
	respBytes, code, err := us.doPostEmptyHeaders(us.buildURL(activities.VersionUrl, activities.ActivitiesUrl), entityData)
	if err != nil {
		return entity, code, err
	}
	return entity, code, json.Unmarshal(respBytes, &entity)
}

// GetActivitiesList retrieves a list of activities based on the provided query parameters.
// It constructs the request URL using the base activities URL and optional query parameters (if provided).
// Returns an `ActivitiesList` containing the list of activities and any error encountered during the request or unmarshalling.
func (us *Uspacy) GetActivitiesList(params url.Values) (entities activities.ActivitiesList, err error) {
	url := us.buildURL(activities.VersionUrl, activities.ActivitiesUrl)
	if len(params) != 0 {
		url = url + "?" + params.Encode()
	}
	body, err := us.doGetEmptyHeaders(url)
	if err != nil {
		return entities, err
	}
	return entities, json.Unmarshal(body, &entities)
}

// GetActivity retrieves details of a specific activity based on its entity ID and optional query parameters.
// The URL is constructed by formatting the activity URL with the given entity ID and appending query parameters if provided.
// Returns the requested `Activity` object and any error encountered during the request or unmarshalling.
func (us *Uspacy) GetActivity(entityId int64, params url.Values) (entity activities.Activity, err error) {
	url := us.buildURL(activities.VersionUrl, fmt.Sprintf(activities.ActivityUrl, strconv.FormatInt(entityId, 10)))
	if len(params) != 0 {
		url = url + "?" + params.Encode()
	}
	body, err := us.doGetEmptyHeaders(url)
	if err != nil {
		return entity, err
	}
	return entity, json.Unmarshal(body, &entity)
}

// PatchActivity updates an existing activity identified by the entity ID with the provided entity data.
// It sends a PATCH request to the constructed URL.
// The function does not return any object, only an error if the request fails or if an issue occurs during the operation.
func (us *Uspacy) PatchActivity(entityId int64, entityData map[string]interface{}) error {
	url := us.buildURL(activities.VersionUrl, fmt.Sprintf(activities.ActivityUrl, strconv.FormatInt(entityId, 10)))
	_, err := us.doPatchEmptyHeaders(url, entityData)
	if err != nil {
		return err
	}
	return nil
}

// DeleteActivity deletes an existing activity identified by the entity ID.
// It sends a DELETE request to the constructed URL and returns the HTTP status code and any error encountered during the request.
// The HTTP status code can be used to check if the deletion was successful.
func (us *Uspacy) DeleteActivity(entityId int64) (int, error) {
	url := us.buildURL(activities.VersionUrl, fmt.Sprintf(activities.ActivityUrl, strconv.FormatInt(entityId, 10)))
	return us.doDeleteEmptyHeaders(url)
}
