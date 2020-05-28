package model

// Order is an order of an cart that can be placed.
type Order struct {
	ID string

	Hash      []byte
	Cart      *Cart
	CartID    string
	Buyer     Address
	Recipient Address
	Coupons   []*Coupon
	Price     int // in cents
	Positions []Position
	Locked    bool
}
