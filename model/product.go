package model

// Product is a product that can be ordered.
type Product struct {
	ID string

	Name  string
	Price int // in cents
}
