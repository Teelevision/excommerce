package controller

import (
	"context"
	"errors"

	"github.com/Teelevision/excommerce/authentication"
	"github.com/Teelevision/excommerce/model"
	"github.com/Teelevision/excommerce/persistence"
)

// Cart is the controller that handles carts.
type Cart struct {
	CartRepository persistence.CartRepository
}

// CreateAndGet creates the given cart. ErrConflict is returned if a cart with
// the same id already exists. The cart is returned with all prices already
// calculated.
func (c *Cart) CreateAndGet(ctx context.Context, cart *model.Cart) (*model.Cart, error) {
	err := c.CartRepository.CreateCart(ctx,
		authentication.AuthenticatedUser(ctx).ID,
		cart.ID,
		convertCartPositions(cart.Positions),
	)
	switch {
	case errors.Is(err, persistence.ErrConflict):
		return nil, ErrConflict
	case err == nil:
		return &model.Cart{
			ID:        cart.ID,
			Positions: calculatePositionPrices(cart.Positions),
		}, nil
	default:
		panic(err)
	}
}

// UpdateAndGet updates the given cart. ErrNotFound is returned if the cart with
// the same id does not exist. ErrForbidden is returned if the cart exists, but
// updating it is not allowed for the current user. The cart is returned with
// all prices already calculated.
func (c *Cart) UpdateAndGet(ctx context.Context, cart *model.Cart) (*model.Cart, error) {
	err := c.CartRepository.UpdateCartOfUser(ctx,
		authentication.AuthenticatedUser(ctx).ID,
		cart.ID,
		convertCartPositions(cart.Positions),
	)
	switch {
	case errors.Is(err, persistence.ErrNotFound):
		return nil, ErrNotFound
	case errors.Is(err, persistence.ErrNotOwnedByUser):
		return nil, ErrForbidden
	case err == nil:
		return &model.Cart{
			ID:        cart.ID,
			Positions: calculatePositionPrices(cart.Positions),
		}, nil
	default:
		panic(err)
	}
}

func calculatePositionPrices(positions []model.Position) []model.Position {
	result := make([]model.Position, len(positions))
	for i, position := range positions {
		position.Price = position.Quantity * position.Product.Price
		result[i] = position
	}
	return result
}

func convertCartPositions(cartPositions []model.Position) (positions []struct {
	ProductID string
	Quantity  int
	Price     int // in cents
}) {
	for _, position := range cartPositions {
		positions = append(positions, struct {
			ProductID string
			Quantity  int
			Price     int
		}{
			ProductID: position.ProductID,
			Quantity:  position.Quantity,
			Price:     position.Price,
		})
	}
	return
}
