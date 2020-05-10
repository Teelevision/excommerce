package persistence

import (
	"context"

	"github.com/Teelevision/excommerce/model"
)

// UserRepository stores and loads users. It is safe for concurrent use.
type UserRepository interface {
	// CreateUser creates a user with the given id, name and password. Id must
	// be unique. Name must be unique. ErrConflict is returned otherwise. The
	// password is stored as a hash and can never be retrieved again.
	CreateUser(ctx context.Context, id, name, password string) error

	// FindUserByNameAndPassword finds the user by the given name and password.
	// As names are unique the result is unambiguous. ErrNotFound is returned if
	// no user matches the set of name and password.
	FindUserByNameAndPassword(ctx context.Context, name, password string) (*model.User, error)

	// FindUserByIDAndPassword finds the user by the given id and password. As
	// ids are unique the result is unambiguous. ErrNotFound is returned if no
	// user matches the set of id and password.
	FindUserByIDAndPassword(ctx context.Context, id, password string) (*model.User, error)
}

// ProductRepository stores and loads products. It is safe for concurrent use.
type ProductRepository interface {
	// CreateProduct creates a product with the given id, name and price. Id
	// must be unique. ErrConflict is returned otherwise. The price is in cents.
	CreateProduct(ctx context.Context, id, name string, price int) error
	// FindAllProducts returns all stored products.
	FindAllProducts(context.Context) ([]*model.Product, error)
	// FindProduct returns the product with the given id. ErrNotFound is
	// returned if there is no product with the id.
	FindProduct(ctx context.Context, id string) (*model.Product, error)
}

// CartRepository stores and loads carts and their positions. It is safe for
// concurrent use.
type CartRepository interface {
	// CreateCart creates a cart for the given user with the given id and
	// positions. Id must be unique. ErrConflict is returned otherwise.
	CreateCart(ctx context.Context, userID, id string, positions []struct {
		ProductID string // TODO: let's see how bad of an idea inlining this really is
		Quantity  int
		Price     int // in cents
	}) error
	// UpdateCartOfUser updates a cart of the given user with new positions. Any
	// existing positions are replaced. ErrNotFound is returned if the cart does
	// not exist. ErrNotOwnedByUser is returned if the cart exists but it's not
	// owned by the given user.
	UpdateCartOfUser(ctx context.Context, userID, id string, positions []struct {
		ProductID string
		Quantity  int
		Price     int // in cents
	}) error
	// FindAllUnlockedCartsOfUser returns all stored carts and their positions
	// of the given user.
	FindAllUnlockedCartsOfUser(ctx context.Context, userID string) ([]*model.Cart, error)
}
