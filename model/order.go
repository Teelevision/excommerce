package model

// Order is an order of an cart that can be placed.
type Order struct {
	ID string

	Cart        *Cart
	CartID      string
	CartVersion int
	Buyer       Address
	Recipient   Address
	Coupons     []string
	Locked      bool
}
