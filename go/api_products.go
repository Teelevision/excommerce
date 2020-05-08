/*
 * ExCommerce
 *
 * ExCommerce is an example commerce system.
 *
 * API version: beta
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

var _ Router = (*ProductsAPI)(nil)

// A ProductsAPI binds http requests to an api service and writes the service results to the http response
type ProductsAPI struct {
	service ProductsAPIServicer
}

// Routes returns all of the api route for the ProductsApiController
func (c *ProductsAPI) Routes() Routes {
	return Routes{
		{
			"GetAllProducts",
			strings.ToUpper("Get"),
			"/beta/products",
			c.GetAllProducts,
		},
		{
			"StoreCouponForProduct",
			strings.ToUpper("Put"),
			"/beta/products/{productId}/coupon/{couponCode}",
			c.StoreCouponForProduct,
		},
	}
}

// GetAllProducts - Get all products
func (c *ProductsAPI) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	result, err := c.service.GetAllProducts()
	if err != nil {
		w.WriteHeader(500)
		return
	}

	EncodeJSONResponse(result, nil, w)
}

// StoreCouponForProduct - Create product coupon
func (c *ProductsAPI) StoreCouponForProduct(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	productID := params["productId"]
	couponCode := params["couponCode"]
	coupon := &Coupon{}
	if err := json.NewDecoder(r.Body).Decode(&coupon); err != nil {
		w.WriteHeader(500)
		return
	}

	result, err := c.service.StoreCouponForProduct(productID, couponCode, *coupon)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	EncodeJSONResponse(result, nil, w)
}
