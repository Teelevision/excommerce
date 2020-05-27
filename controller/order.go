package controller

import (
	"context"
	"errors"

	"github.com/Teelevision/excommerce/authentication"
	"github.com/Teelevision/excommerce/model"
	"github.com/Teelevision/excommerce/persistence"
	"github.com/google/uuid"
)

// Order is the controller that handles orders.
type Order struct {
	OrderRepository persistence.OrderRepository
}

// CreateAndGet creates the given order. The order is returned with a unique id.
func (c *Order) CreateAndGet(ctx context.Context, order *model.Order) (*model.Order, error) {
	// create id
	uuid, err := uuid.NewRandom()
	if err != nil {
		panic(err)
	}
	id := uuid.String()

	// cart version
	cartVersion := 0 // TODO
	// TODO: We actually need a hash here, that includes the cart price, product
	// prices, the applied coupons and any other discount.

	// coupon codes
	couponCodes := make([]string, len(order.Coupons))
	for i, coupon := range order.Coupons {
		couponCodes[i] = coupon.Code
	}

	err = c.OrderRepository.CreateOrder(ctx,
		authentication.AuthenticatedUser(ctx).ID,
		id,
		persistence.OrderAttributes{
			CartID:      order.CartID,
			CartVersion: cartVersion,
			Buyer:       persistence.OrderAddress(order.Buyer),
			Recipient:   persistence.OrderAddress(order.Recipient),
			Coupons:     couponCodes,
		},
	)
	switch {
	case errors.Is(err, persistence.ErrConflict):
		panic(err)
	case err == nil:
		positions := generateOrderPositions(order.Cart.Positions, order.Coupons)
		positions = calculatePositionPrices(positions)
		return &model.Order{
			ID:          id,
			Buyer:       order.Buyer,
			Recipient:   order.Recipient,
			Cart:        order.Cart,
			CartID:      order.CartID,
			CartVersion: cartVersion,
			Coupons:     order.Coupons,
			Positions:   positions,
			Price:       calculatePositionSum(positions),
		}, nil
	default:
		panic(err)
	}
}

func generateOrderPositions(positions []model.Position, coupons []*model.Coupon) []model.Position {
	// get best coupons
	productCoupons := make(map[string]*model.Coupon, len(coupons))
	for _, coupon := range coupons {
		existingCoupon, ok := productCoupons[coupon.ProductID]
		if !ok || coupon.Discount > existingCoupon.Discount {
			productCoupons[coupon.ProductID] = coupon
		}
	}
	out := make([]model.Position, 0, len(positions)+len(productCoupons))
	for _, position := range positions {
		out = append(out, position)
		coupon, ok := productCoupons[position.ProductID]
		if !ok {
			continue
		}
		// add position for coupon
		price := -coupon.Discount * position.Price / 100
		out = append(out, model.Position{
			Quantity: 1,
			Price:    price,
			Product: &model.Product{
				Name:  coupon.Name,
				Price: price,
			},
		})
	}
	return out
}

func calculatePositionSum(positions []model.Position) (sum int) {
	for _, position := range positions {
		sum += position.Price
	}
	return
}
