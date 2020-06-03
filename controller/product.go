package controller

import (
	"context"
	"errors"
	"time"

	"github.com/Teelevision/excommerce/config"
	"github.com/Teelevision/excommerce/model"
	"github.com/Teelevision/excommerce/persistence"
)

// Product is the controller that handles products.
type Product struct {
	ProductRepository persistence.ProductRepository
	CouponRepository  persistence.CouponRepository
}

// GetAll gets all products.
func (c *Product) GetAll(ctx context.Context) ([]*model.Product, error) {
	products, err := c.ProductRepository.FindAllProducts(ctx)
	switch {
	case err == nil:
		return products, nil
	default:
		panic(err)
	}
}

// Get returns the requested product. ErrNotFound is returned if there is no
// product with the given id.
func (c *Product) Get(ctx context.Context, productID string) (*model.Product, error) {
	if product := getSpecialProduct(productID, 0); product != nil {
		return product, nil
	}
	product, err := c.ProductRepository.FindProduct(ctx, productID)
	switch {
	case errors.Is(err, persistence.ErrNotFound):
		return nil, ErrNotFound
	case err == nil:
		return &model.Product{
			ID:    product.ID,
			Name:  product.Name,
			Price: product.Price,
		}, nil
	default:
		panic(err)
	}
}

// SaveCoupon creates or updates the given coupon. The coupon's code is expected
// to be 6 to 40 runes long, and the coupon's name 1 to 100. The coupon's
// product is expected to exist and the coupon's discount is expected to be
// between 1 and 100. On success the coupon is returned.
func (c *Product) SaveCoupon(ctx context.Context, coupon *model.Coupon) (*model.Coupon, error) {
	if coupon.ExpiresAt.IsZero() {
		coupon.ExpiresAt = time.Now().Add(config.CouponDefaultLifetime)
	}

	err := c.CouponRepository.StoreCoupon(ctx, coupon.Code, coupon.Name, coupon.Product.ID, coupon.Discount, coupon.ExpiresAt)
	switch {
	case err == nil:
		return coupon, nil
	default:
		panic(err)
	}
}

// GetCoupon returns the valid coupon with the given code. ErrNotFound is
// returned if there is no coupon with the given code or it is invalid.
func (c *Product) GetCoupon(ctx context.Context, code string) (*model.Coupon, error) {
	coupon, err := c.CouponRepository.FindValidCoupon(ctx, code)
	switch {
	case errors.Is(err, persistence.ErrNotFound):
		return nil, ErrNotFound
	case err == nil:
		return coupon, nil
	default:
		panic(err)
	}
}

func getSpecialProduct(productID string, price int) *model.Product {
	switch productID {
	case "0de17a66-ea59-4032-9383-2603c6c77d25": // set of 4 pears and 2 bananas
		return &model.Product{
			ID:    "0de17a66-ea59-4032-9383-2603c6c77d25",
			Name:  "Set of 4 pears and 2 bananas (30% off)",
			Price: 444,
		}
	default:
		return nil
	}
}
