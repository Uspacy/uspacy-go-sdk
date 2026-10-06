package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/Uspacy/uspacy-go-sdk/v2/crm"
)

// GetProduct returns list of products
func (us *Uspacy) GetProduct(ctx context.Context, id string) (call crm.Products, err error) {
	responseBody, err := us.doGetEmptyHeaders(ctx, us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.ProductsUrl, id)))
	if err != nil {
		return call, err
	}
	return call, json.Unmarshal(responseBody, &call)
}

// GetEntityProductList returns product list for entity
func (us *Uspacy) GetEntityProductList(ctx context.Context, entityType string, entityID int64) (productList crm.EntityProductList, err error) {
	params := url.Values{
		"entity_type": []string{entityType},
		"entity_id":   []string{strconv.FormatInt(entityID, 10)},
	}
	responseBody, err := us.doGetEmptyHeaders(ctx, us.buildURL(crm.VersionUrl, crm.EntityProductListsUrl)+"?"+params.Encode())
	if err != nil {
		return productList, err
	}
	return productList, json.Unmarshal(responseBody, &productList)
}

// CreateEntityListProduct creates a product list for entity
func (us *Uspacy) CreateEntityListProduct(ctx context.Context, productData crm.CreateEntityListProductRequest, opts ...RequestOption) (product crm.EntityListProduct, err error) {
	responseBody, _, err := us.doPost(ctx, us.buildURL(crm.VersionUrl, crm.ListProductsUrl), productData, opts...)
	if err != nil {
		return product, err
	}
	return product, json.Unmarshal(responseBody, &product)
}

// DeleteEntityListProduct deletes a product list item
func (us *Uspacy) DeleteEntityListProduct(ctx context.Context, id int) (statusCode int, err error) {
	return us.doDeleteEmptyHeaders(ctx, us.buildURL(crm.VersionUrl, crm.ListProductsUrl, strconv.Itoa(id)), nil)
}

// CreateProduct returns list of products
func (us *Uspacy) CreateProduct(ctx context.Context, productData map[string]any, opts ...RequestOption) (product crm.Products, err error) {
	responseBody, _, err := us.doPost(ctx, us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.ProductsUrl, "")), productData, opts...)
	if err != nil {
		return product, err
	}
	return product, json.Unmarshal(responseBody, &product)
}
