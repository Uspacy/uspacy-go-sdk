package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"uspacy-go-sdk/crm"
)

// GetFields returns Fields struct for a given type of field
func (us *Uspacy) GetFields(field string) (crm.Fields, error) {

	var fields crm.Fields

	req, err := us.generateRequest(mainHost+fmt.Sprintf(crm.ListFields, field), http.MethodGet, emptyHeaders, nil)
	if err != nil {
		return fields, err
	}

	body, err := us.doRaw(req)
	if err != nil {
		return fields, err
	}

	if err = json.Unmarshal(body, &fields); err != nil {
		return fields, err
	}

	return fields, nil
}
