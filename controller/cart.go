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
	CartRepository    persistence.CartRepository
	ProductRepository persistence.ProductRepository
}

// Get returns the cart with the given id with all prices calculated.
// ErrNotFound is retuned if there is no cart with the id. ErrDeleted is
// returned if the cart did exist but is deleted. ErrForbidden is returned if
// the cart exists, but the current user is not allowed to access it.
func (c *Cart) Get(ctx context.Context, cartID string) (*model.Cart, error) {
	cart, err := c.CartRepository.FindCartOfUser(ctx,
		authentication.AuthenticatedUser(ctx).ID,
		cartID,
	)
	switch {
	case errors.Is(err, persistence.ErrNotFound):
		return nil, ErrNotFound
	case errors.Is(err, persistence.ErrDeleted):
		return nil, ErrDeleted
	case errors.Is(err, persistence.ErrNotOwnedByUser):
		return nil, ErrForbidden
	case err == nil:
		// load products
		c.loadProducts(ctx, cart)
		cart.Positions = calculatePositionPrices(cart.Positions)
		return cart, nil
	default:
		panic(err)
	}
}

// GetAllUnlocked returns all unlocked carts of the current user.
func (c *Cart) GetAllUnlocked(ctx context.Context) ([]*model.Cart, error) {
	carts, err := c.CartRepository.FindAllUnlockedCartsOfUser(ctx,
		authentication.AuthenticatedUser(ctx).ID)
	switch {
	case err == nil:
		for _, cart := range carts {
			c.loadProducts(ctx, cart)
			cart.Positions = calculatePositionPrices(cart.Positions)
		}
		return carts, nil
	default:
		panic(err)
	}
}

// CreateAndGet creates the given cart. ErrConflict is returned if a cart with
// the same id already exists or existed. The cart is returned with all prices
// already calculated.
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
		cart.Positions = calculatePositionPrices(cart.Positions)
		return cart, nil
	default:
		panic(err)
	}
}

// UpdateAndGet updates the given cart. ErrNotFound is returned if the cart with
// the same id does not exist. ErrDeleted is returned if the cart did exist but
// is deleted. ErrForbidden is returned if the cart exists, but updating it is
// not allowed for the current user. The cart is returned with all prices
// already calculated.
func (c *Cart) UpdateAndGet(ctx context.Context, cart *model.Cart) (*model.Cart, error) {
	err := c.CartRepository.UpdateCartOfUser(ctx,
		authentication.AuthenticatedUser(ctx).ID,
		cart.ID,
		convertCartPositions(cart.Positions),
	)
	switch {
	case errors.Is(err, persistence.ErrNotFound):
		return nil, ErrNotFound
	case errors.Is(err, persistence.ErrDeleted):
		return nil, ErrDeleted
	case errors.Is(err, persistence.ErrNotOwnedByUser):
		return nil, ErrForbidden
	case err == nil:
		cart.Positions = calculatePositionPrices(cart.Positions)
		return cart, nil
	default:
		panic(err)
	}
}

// Delete deletes the cart with the given id. ErrNotFound is retuned if there is
// no cart with the id. ErrDeleted is returned if the cart did exist but is
// deleted. ErrForbidden is returned if the cart exists, but the current user is
// not allowed to delete it.
func (c *Cart) Delete(ctx context.Context, cartID string) error {
	err := c.CartRepository.DeleteCartOfUser(ctx,
		authentication.AuthenticatedUser(ctx).ID,
		cartID,
	)
	switch {
	case errors.Is(err, persistence.ErrNotFound):
		return ErrNotFound
	case errors.Is(err, persistence.ErrDeleted):
		return ErrDeleted
	case errors.Is(err, persistence.ErrNotOwnedByUser):
		return ErrForbidden
	case err == nil:
		return nil
	default:
		panic(err)
	}
}

func (c *Cart) loadProducts(ctx context.Context, cart *model.Cart) {
	for i, position := range cart.Positions {
		product, err := c.ProductRepository.FindProduct(ctx, position.ProductID)
		switch {
		case errors.Is(err, persistence.ErrNotFound):
			cart.Positions[i].ProductID = ""
			cart.Positions[i].Product = &model.Product{Name: "Product not available anymore."}
		case err == nil:
			cart.Positions[i].Product = product
		default:
			panic(err)
		}
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
