package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	types "uspacy-go-sdk"
)

// GetFields returns Fields struct for a given type of field
func (us *Uspacy) GetFields(field string) (types.Fields, error) {

	var fields types.Fields
	uri := mainHost + fmt.Sprintf(crmListFields, field)

	req, err := us.generateRequest(uri, http.MethodGet, emptyHeaders, nil)
	if err != nil {
		return fields, err
	}

	res, err := us.client.Do(req)
	if err != nil {
		return fields, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fields, err
	}

	if !handleStatusCode(res.StatusCode) {
		log.Fatalf("error occured while trying to GetFields\nbody - %s\ncode - %v", string(body), res.StatusCode)
		return fields, err

	}

	if err = json.Unmarshal(body, &fields); err != nil {
		return fields, err
	}

	return fields, nil
}
