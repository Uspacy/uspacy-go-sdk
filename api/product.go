package api

import (
	"encoding/json"
	"fmt"

	"github.com/Uspacy/uspacy-go-sdk/crm"
)

// GetProduct returns list of products
func (us *Uspacy) GetProduct(id string) (call crm.Products, err error) {
	responseBody, err := us.doGetEmptyHeaders(us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.ProductsUrl, id)))
	if err != nil {
		return call, err
	}
	return call, json.Unmarshal(responseBody, &call)
}

// CreateProduct returns list of products
func (us *Uspacy) CreateProduct(productData map[string]interface{}) (product crm.Products, err error) {
	responseBody, _, err := us.doPostEmptyHeaders(us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.ProductsUrl, "")), productData)
	if err != nil {
		return product, err
	}
	return product, json.Unmarshal(responseBody, &product)
}
