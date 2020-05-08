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

var _ Router = (*CartsAPI)(nil)

// A CartsAPI binds http requests to an api service and writes the service results to the http response
type CartsAPI struct {
	service CartsAPIServicer
}

// Routes returns all of the api route for the CartsApiController
func (c *CartsAPI) Routes() Routes {
	return Routes{
		{
			"DeleteCart",
			strings.ToUpper("Delete"),
			"/beta/carts/{cartId}",
			c.DeleteCart,
		},
		{
			"GetAllCarts",
			strings.ToUpper("Get"),
			"/beta/carts",
			c.GetAllCarts,
		},
		{
			"GetCart",
			strings.ToUpper("Get"),
			"/beta/carts/{cartId}",
			c.GetCart,
		},
		{
			"StoreCart",
			strings.ToUpper("Put"),
			"/beta/carts/{cartId}",
			c.StoreCart,
		},
	}
}

// DeleteCart - Delete a cart
func (c *CartsAPI) DeleteCart(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	cartID := params["cartId"]
	result, err := c.service.DeleteCart(cartID)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	EncodeJSONResponse(result, nil, w)
}

// GetAllCarts - Get all carts
func (c *CartsAPI) GetAllCarts(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	locked := query.Get("locked") == "true"
	result, err := c.service.GetAllCarts(locked)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	EncodeJSONResponse(result, nil, w)
}

// GetCart - Get a cart
func (c *CartsAPI) GetCart(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	cartID := params["cartId"]
	result, err := c.service.GetCart(cartID)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	EncodeJSONResponse(result, nil, w)
}

// StoreCart - Store a cart
func (c *CartsAPI) StoreCart(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	cartID := params["cartId"]
	cart := &Cart{}
	if err := json.NewDecoder(r.Body).Decode(&cart); err != nil {
		w.WriteHeader(500)
		return
	}

	result, err := c.service.StoreCart(cartID, *cart)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	EncodeJSONResponse(result, nil, w)
}
