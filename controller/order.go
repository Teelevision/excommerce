package controller

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"

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

	// prepare positions
	positions := calculateCartPositionPrices(order.Cart.Positions)
	positions = generateOrderPositions(positions, order.Coupons)

	// hash
	hash := hashPositions(positions)

	// coupon codes
	couponCodes := make([]string, len(order.Coupons))
	for i, coupon := range order.Coupons {
		couponCodes[i] = coupon.Code
	}

	err = c.OrderRepository.CreateOrder(ctx,
		authentication.AuthenticatedUser(ctx).ID,
		id,
		persistence.OrderAttributes{
			Hash:      hash,
			CartID:    order.CartID,
			Buyer:     persistence.OrderAddress(order.Buyer),
			Recipient: persistence.OrderAddress(order.Recipient),
			Coupons:   couponCodes,
		},
	)
	switch {
	case errors.Is(err, persistence.ErrConflict):
		panic(err)
	case err == nil:
		return &model.Order{
			ID:        id,
			Hash:      hash,
			Buyer:     order.Buyer,
			Recipient: order.Recipient,
			Cart:      order.Cart,
			CartID:    order.CartID,
			Coupons:   order.Coupons,
			Positions: positions,
			Price:     calculatePositionSum(positions),
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
			Quantity:   1,
			Price:      price,
			Coupon:     coupon,
			CouponCode: coupon.Code,
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

func hashPositions(positions []model.Position) []byte {
	entries := make(sort.StringSlice, len(positions))
	for i, position := range positions {
		buf := new(bytes.Buffer)
		fmt.Fprintf(buf, "%d,%d,", position.Quantity, position.Price)
		switch {
		case position.Product != nil:
			fmt.Fprintf(buf, "product:%s", position.Product.ID)
		case position.Coupon != nil:
			fmt.Fprintf(buf, "coupon:%s,%d,%q", position.Coupon.ProductID, position.Coupon.Discount, position.Coupon.Code)
		default:
			panic("position has no product and no coupon")
		}
		entries[i] = buf.String()
	}
	entries.Sort()
	base := strings.Join(entries, "\n")
	// TODO: choose hash algorithm
	return []byte(base)
}
