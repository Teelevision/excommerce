package controller

import (
	"context"
	"errors"

	"github.com/Teelevision/excommerce/authentication"
	"github.com/Teelevision/excommerce/model"
	"github.com/Teelevision/excommerce/persistence"
)

// StoreCartController is the controller that stores carts.
type StoreCartController struct {
	CartRepository persistence.CartRepository
}

// Create creates the given cart. ErrConflict is returned if a cart with the
// same id already exists.
func (c *StoreCartController) Create(ctx context.Context, cart *model.Cart) error {
	err := c.CartRepository.CreateCart(ctx,
		authentication.AuthenticatedUser(ctx).ID,
		cart.ID,
		convertCartPositions(cart.Positions),
	)
	switch {
	case errors.Is(err, persistence.ErrConflict):
		return ErrConflict
	case err == nil:
		return nil
	default:
		panic(err)
	}
}

// Update updates the given cart. ErrNotFound is returned if the cart with the
// same id does not exist. ErrForbidden is returned if the cart exists, but
// updating it is not allowed for the current user.
func (c *StoreCartController) Update(ctx context.Context, cart *model.Cart) error {
	err := c.CartRepository.UpdateCartOfUser(ctx,
		authentication.AuthenticatedUser(ctx).ID,
		cart.ID,
		convertCartPositions(cart.Positions),
	)
	switch {
	case errors.Is(err, persistence.ErrNotFound):
		return ErrNotFound
	case errors.Is(err, persistence.ErrNotOwnedByUser):
		return ErrForbidden
	case err == nil:
		return nil
	default:
		panic(err)
	}
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
