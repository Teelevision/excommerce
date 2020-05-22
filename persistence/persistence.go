package persistence

import (
	"context"
	"time"

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
	// Positions maps product ids to quantity.
	CreateCart(ctx context.Context, userID, id string, positions map[string]int) error
	// UpdateCartOfUser updates a cart of the given user with new positions. Any
	// existing positions are replaced. ErrNotFound is returned if the cart does
	// not exist. ErrDeleted is returned if the cart did exist but is deleted.
	// ErrNotOwnedByUser is returned if the cart exists but it's not owned by
	// the given user.
	UpdateCartOfUser(ctx context.Context, userID, id string, positions map[string]int) error
	// FindAllUnlockedCartsOfUser returns all stored carts and their positions
	// of the given user.
	FindAllUnlockedCartsOfUser(ctx context.Context, userID string) ([]*model.Cart, error)
	// FindCartOfUser returns the cart of the given user with the given cart id.
	// ErrNotFound is returned if there is no cart with the id. ErrDeleted is
	// returned if the cart did exist but is deleted. ErrNotOwnedByUser is
	// returned if the cart exists but it's not owned by the given user.
	FindCartOfUser(ctx context.Context, userID, id string) (*model.Cart, error)
	// DeleteCartOfUser deletes the cart of the given user with the given cart
	// id. ErrNotFound is returned if there is no cart with the id. ErrDeleted
	// is returned if the cart did exist but is deleted. ErrNotOwnedByUser is
	// returned if the cart exists but it's not owned by the given user.
	DeleteCartOfUser(ctx context.Context, userID, id string) error
}

// CouponRepository stores and loads coupons. It is safe for concurrent use.
type CouponRepository interface {
	// StoreCoupon stores a coupon with the given code, name, product id,
	// discount in percent and expires at time. If a coupon with the same code
	// was previously stored it is overwritten.
	StoreCoupon(ctx context.Context, code, name, productID string, discount int, expiresAt time.Time) error
	// FindValidCoupon returns the coupon with the given code that is not
	// expired. ErrNotFound is returned if there is no coupon with the code or
	// the coupon is expired.
	FindValidCoupon(ctx context.Context, code string) (*model.Coupon, error)
}

// OrderRepository stores and loads orders. It is safe for concurrent use.
type OrderRepository interface {
	// CreateOrder creates an order for the given user with the given id and
	// attributes. Id must be unique. ErrConflict is returned otherwise.
	CreateOrder(ctx context.Context, userID, id string, attributes OrderAttributes) error
	// FindOrderOfUser returns the order of the given user with the given id.
	// ErrNotFound is returned if there is no order with the id. ErrDeleted is
	// returned if the order did exist but is deleted. ErrNotOwnedByUser is
	// returned if the order exists but it's not owned by the given user.
	FindOrderOfUser(ctx context.Context, userID, id string) (*model.Order, error)
	// DeleteOrderOfUser deletes the order of the given user with the given id.
	// ErrNotFound is returned if there is no order with the id. ErrDeleted is
	// returned if the order did exist but is deleted. ErrNotOwnedByUser is
	// returned if the order exists but it's not owned by the given user.
	// ErrLocked is returned if the order is owned by the given user, but is
	// locked.
	DeleteOrderOfUser(ctx context.Context, userID, id string) error
	// LockOrderOfUser locks the order of the given user with the given id.
	// ErrNotFound is returned if there is no order with the id. ErrDeleted is
	// returned if the order did exist but is deleted. ErrNotOwnedByUser is
	// returned if the order exists but it's not owned by the given user.
	// ErrLocked is returned if the order is owned by the given user, but is
	// locked.
	LockOrderOfUser(ctx context.Context, userID, id string) error
}

// OrderAttributes are common attributes of an order.
type OrderAttributes struct {
	CartID      string
	CartVersion int
	Buyer       OrderAddress
	Recipient   OrderAddress
	Coupons     []string
}

// OrderAddress is an address used in orders.
type OrderAddress struct {
	Name       string
	Country    string
	PostalCode string
	City       string
	Street     string
}
