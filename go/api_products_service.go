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
	"errors"
)

// ProductsApiService is a service that implents the logic for the ProductsApiServicer
// This service should implement the business logic for every endpoint for the ProductsApi API. 
// Include any external packages or services that will be required by this service.
type ProductsApiService struct {
}

// NewProductsApiService creates a default api service
func NewProductsApiService() ProductsApiServicer {
	return &ProductsApiService{}
}

// GetAllProducts - Get all products
func (s *ProductsApiService) GetAllProducts() (interface{}, error) {
	// TODO - update GetAllProducts with the required logic for this service method.
	// Add api_products_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.
	return nil, errors.New("service method 'GetAllProducts' not implemented")
}

// StoreCouponForProduct - Create product coupon
func (s *ProductsApiService) StoreCouponForProduct(productId string, couponCode string, coupon Coupon) (interface{}, error) {
	// TODO - update StoreCouponForProduct with the required logic for this service method.
	// Add api_products_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.
	return nil, errors.New("service method 'StoreCouponForProduct' not implemented")
}
