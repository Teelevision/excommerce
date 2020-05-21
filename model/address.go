package model

// Address is an address for billing or shipping.
type Address struct {
	Name       string
	Country    string
	PostalCode string
	City       string
	Street     string
}
