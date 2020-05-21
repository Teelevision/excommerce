package model

import "time"

// Coupon is a coupon that gives a discount on a specific product.
type Coupon struct {
	Code string // like an id

	Product   *Product
	ProductID string
	Name      string
	Discount  int // in percent
	ExpiresAt time.Time
}
