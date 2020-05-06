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
	"net/http"
)


// CartsApiRouter defines the required methods for binding the api requests to a responses for the CartsApi
// The CartsApiRouter implementation should parse necessary information from the http request, 
// pass the data to a CartsApiServicer to perform the required actions, then write the service results to the http response.
type CartsApiRouter interface { 
	DeleteCart(http.ResponseWriter, *http.Request)
	GetAllCarts(http.ResponseWriter, *http.Request)
	GetCart(http.ResponseWriter, *http.Request)
	StoreCart(http.ResponseWriter, *http.Request)
}
// OrdersApiRouter defines the required methods for binding the api requests to a responses for the OrdersApi
// The OrdersApiRouter implementation should parse necessary information from the http request, 
// pass the data to a OrdersApiServicer to perform the required actions, then write the service results to the http response.
type OrdersApiRouter interface { 
	CreateOrderFromCart(http.ResponseWriter, *http.Request)
	PlaceOrder(http.ResponseWriter, *http.Request)
}
// ProductsApiRouter defines the required methods for binding the api requests to a responses for the ProductsApi
// The ProductsApiRouter implementation should parse necessary information from the http request, 
// pass the data to a ProductsApiServicer to perform the required actions, then write the service results to the http response.
type ProductsApiRouter interface { 
	GetAllProducts(http.ResponseWriter, *http.Request)
	StoreCouponForProduct(http.ResponseWriter, *http.Request)
}
// UsersApiRouter defines the required methods for binding the api requests to a responses for the UsersApi
// The UsersApiRouter implementation should parse necessary information from the http request, 
// pass the data to a UsersApiServicer to perform the required actions, then write the service results to the http response.
type UsersApiRouter interface { 
	Login(http.ResponseWriter, *http.Request)
	Register(http.ResponseWriter, *http.Request)
}


// CartsApiServicer defines the api actions for the CartsApi service
// This interface intended to stay up to date with the openapi yaml used to generate it, 
// while the service implementation can ignored with the .openapi-generator-ignore file 
// and updated with the logic required for the API.
type CartsApiServicer interface { 
	DeleteCart(string) (interface{}, error)
	GetAllCarts(bool) (interface{}, error)
	GetCart(string) (interface{}, error)
	StoreCart(string, Cart) (interface{}, error)
}


// OrdersApiServicer defines the api actions for the OrdersApi service
// This interface intended to stay up to date with the openapi yaml used to generate it, 
// while the service implementation can ignored with the .openapi-generator-ignore file 
// and updated with the logic required for the API.
type OrdersApiServicer interface { 
	CreateOrderFromCart(string, Order) (interface{}, error)
	PlaceOrder(string) (interface{}, error)
}


// ProductsApiServicer defines the api actions for the ProductsApi service
// This interface intended to stay up to date with the openapi yaml used to generate it, 
// while the service implementation can ignored with the .openapi-generator-ignore file 
// and updated with the logic required for the API.
type ProductsApiServicer interface { 
	GetAllProducts() (interface{}, error)
	StoreCouponForProduct(string, string, Coupon) (interface{}, error)
}


// UsersApiServicer defines the api actions for the UsersApi service
// This interface intended to stay up to date with the openapi yaml used to generate it, 
// while the service implementation can ignored with the .openapi-generator-ignore file 
// and updated with the logic required for the API.
type UsersApiServicer interface { 
	Login(LoginForm) (interface{}, error)
	Register(User) (interface{}, error)
}
