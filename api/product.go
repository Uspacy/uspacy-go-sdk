package api

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

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

func (us *Uspacy) GetEntityProductList(entityType string, entityID int64) (productList crm.EntityProductList, err error) {
	params := url.Values{
		"entity_type": []string{entityType},
		"entity_id":   []string{strconv.FormatInt(entityID, 10)},
	}
	responseBody, err := us.doGetEmptyHeaders(us.buildURL(crm.VersionUrl, crm.EntityProductListsUrl) + "?" + params.Encode())
	if err != nil {
		return productList, err
	}
	return productList, json.Unmarshal(responseBody, &productList)
}

func (us *Uspacy) CreateEntityListProduct(productData crm.CreateEntityListProductRequest, headers ...map[string]string) (product crm.EntityListProduct, err error) {
	responseBody, _, err := us.doPost(us.buildURL(crm.VersionUrl, crm.ListProductsUrl), productData, headers...)
	if err != nil {
		return product, err
	}
	return product, json.Unmarshal(responseBody, &product)
}

func (us *Uspacy) DeleteEntityListProduct(id int) (statusCode int, err error) {
	return us.doDeleteEmptyHeaders(us.buildURL(crm.VersionUrl, crm.ListProductsUrl, strconv.Itoa(id)), nil)
}

// CreateProduct returns list of products
func (us *Uspacy) CreateProduct(productData map[string]any, headers ...map[string]string) (product crm.Products, err error) {
	responseBody, _, err := us.doPost(us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.ProductsUrl, "")), productData, headers...)
	if err != nil {
		return product, err
	}
	return product, json.Unmarshal(responseBody, &product)
}
