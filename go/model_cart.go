/*
 * ExCommerce
 *
 * ExCommerce is an example commerce system.
 *
 * API version: beta
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

// Cart - A cart containing products.
type Cart struct {

	// The UUID of the cart.
	Id string `json:"id"`

	Positions []Position `json:"positions"`

	// Whether the cart is locked.
	Locked bool `json:"locked,omitempty"`
}
