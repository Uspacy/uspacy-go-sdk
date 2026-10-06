package api

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/Uspacy/uspacy-go-sdk/v2/crm"
)

// GetProduct fetches static/products/{id}; an empty id lists all products. The
// response is decoded as crm.Products.
func (us *Uspacy) GetProduct(ctx context.Context, id string, opts ...RequestOption) (crm.Products, error) {
	responseBody, err := us.doGet(ctx, us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.ProductsUrl, id)), opts...)
	return decodeJSON[crm.Products](responseBody, err)
}

// GetEntityProductList returns the product list of a CRM record.
func (us *Uspacy) GetEntityProductList(ctx context.Context, entityType string, entityID int64, opts ...RequestOption) (crm.EntityProductList, error) {
	params := url.Values{
		"entity_type": []string{entityType},
		"entity_id":   []string{strconv.FormatInt(entityID, 10)},
	}
	responseBody, err := us.doGet(ctx, us.buildURL(crm.VersionUrl, crm.EntityProductListsUrl)+"?"+params.Encode(), opts...)
	return decodeJSON[crm.EntityProductList](responseBody, err)
}

// CreateEntityListProduct adds a product to a CRM record's product list and returns the
// new item.
func (us *Uspacy) CreateEntityListProduct(ctx context.Context, productData crm.CreateEntityListProductRequest, opts ...RequestOption) (crm.EntityListProduct, error) {
	responseBody, _, err := us.doPost(ctx, us.buildURL(crm.VersionUrl, crm.ListProductsUrl), productData, opts...)
	return decodeJSON[crm.EntityListProduct](responseBody, err)
}

// DeleteEntityListProduct removes an item from a CRM record's product list and returns the
// HTTP status code.
func (us *Uspacy) DeleteEntityListProduct(ctx context.Context, id int, opts ...RequestOption) (statusCode int, err error) {
	return us.doDelete(ctx, us.buildURL(crm.VersionUrl, crm.ListProductsUrl, strconv.Itoa(id)), nil, opts...)
}

// CreateProduct creates a product. The response is decoded as crm.Products.
func (us *Uspacy) CreateProduct(ctx context.Context, productData map[string]any, opts ...RequestOption) (crm.Products, error) {
	responseBody, _, err := us.doPost(ctx, us.buildURL(crm.VersionUrl, fmt.Sprintf(crm.ProductsUrl, "")), productData, opts...)
	return decodeJSON[crm.Products](responseBody, err)
}
