package model

// Position is a position of a product in a cart or order.
type Position struct {
	Product    *Product
	ProductID  string
	Coupon     *Coupon
	CouponCode string
	Quantity   int
	Price      int // in cents
}
