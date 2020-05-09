package model

// Cart is a cart that contains products.
type Cart struct {
	ID string

	Positions []Position
	Locked    bool
}
