package inmemory

import (
	"context"
	"sync"

	"github.com/Teelevision/excommerce/model"
	"github.com/Teelevision/excommerce/persistence"
	"golang.org/x/crypto/bcrypt"
)

// Adapter is the in-memory persistence adapter. It implements a range of
// repositories. Please use NewAdapter to create a new instance. Adapter is safe
// for concurrent use.
type Adapter struct {
	mx sync.Mutex

	usersByID    map[string]*user
	usersByName  map[string]*user
	productsByID map[string]*product
	cartsByID    map[string]*cart

	bcryptCost int
}

// Option can be used to configure an adapter.
type Option func(*Adapter)

// FastLessSecureHashingForTesting is an option that configures an adapter to
// use a less secure hashing. It is not secure enough to use in production, but
// can speed up tests.
func FastLessSecureHashingForTesting() Option {
	return func(a *Adapter) {
		a.bcryptCost = bcrypt.MinCost
	}
}

// NewAdapter returns a new in-memory adapter.
func NewAdapter(options ...Option) *Adapter {
	a := Adapter{
		usersByID:    make(map[string]*user),
		usersByName:  make(map[string]*user),
		productsByID: make(map[string]*product),
		cartsByID:    make(map[string]*cart),
	}
	for _, option := range options {
		option(&a)
	}
	return &a
}

var _ persistence.UserRepository = (*Adapter)(nil)

type user struct {
	id           string
	name         string
	passwordHash []byte // bcrypt
}

// CreateUser creates a user with the given id, name and password. Id must be
// unique. Name must be unique. ErrConflict is returned otherwise. The password
// is stored as a hash and can never be retrieved again.
func (a *Adapter) CreateUser(_ context.Context, id string, name string, password string) error {
	a.mx.Lock()
	defer a.mx.Unlock()

	// check that id is unique
	if _, ok := a.usersByID[id]; ok {
		return persistence.ErrConflict
	}
	// check that name is unique
	if _, ok := a.usersByName[name]; ok {
		return persistence.ErrConflict
	}

	// hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), a.bcryptCost)
	if err != nil {
		panic(err)
	}

	// save user
	user := user{
		id:           id,
		name:         name,
		passwordHash: hash,
	}
	a.usersByID[id] = &user
	a.usersByName[name] = &user
	return nil
}

// FindUserByNameAndPassword finds the user by the given name and password. As
// names are unique the result is unambiguous. ErrNotFound is returned if no
// user matches the set of name and password.
func (a *Adapter) FindUserByNameAndPassword(_ context.Context, name string, password string) (*model.User, error) {
	a.mx.Lock()
	defer a.mx.Unlock()

	user, ok := a.usersByName[name]
	if !ok {
		return nil, persistence.ErrNotFound
	}
	return checkUserPassword(user, password)
}

// FindUserByIDAndPassword finds the user by the given id and password. As ids
// are unique the result is unambiguous. ErrNotFound is returned if no user
// matches the set of id and password.
func (a *Adapter) FindUserByIDAndPassword(_ context.Context, id string, password string) (*model.User, error) {
	a.mx.Lock()
	defer a.mx.Unlock()

	user, ok := a.usersByID[id]
	if !ok {
		return nil, persistence.ErrNotFound
	}
	return checkUserPassword(user, password)
}

func checkUserPassword(user *user, password string) (*model.User, error) {
	// check password
	if err := bcrypt.CompareHashAndPassword(user.passwordHash, []byte(password)); err != nil {
		return nil, persistence.ErrNotFound
	}
	return &model.User{
		ID:   user.id,
		Name: user.name,
	}, nil
}

var _ persistence.ProductRepository = (*Adapter)(nil)

type product struct {
	name  string
	price int // in cents
}

// CreateProduct creates a product with the given id, name and price. Id must be
// unique. ErrConflict is returned otherwise. The price is in cents.
func (a *Adapter) CreateProduct(_ context.Context, id, name string, price int) error {
	a.mx.Lock()
	defer a.mx.Unlock()

	if _, ok := a.productsByID[id]; ok {
		return persistence.ErrConflict
	}

	a.productsByID[id] = &product{
		name:  name,
		price: price,
	}
	return nil
}

// FindAllProducts returns all stored products.
func (a *Adapter) FindAllProducts(_ context.Context) ([]*model.Product, error) {
	a.mx.Lock()
	defer a.mx.Unlock()

	result := make([]*model.Product, 0, len(a.productsByID))
	for id, product := range a.productsByID {
		result = append(result, &model.Product{
			ID:    id,
			Name:  product.name,
			Price: product.price,
		})
	}
	return result, nil
}

// FindProduct returns the product with the given id. ErrNotFound is returned if
// there is no product with the id.
func (a *Adapter) FindProduct(ctx context.Context, id string) (*model.Product, error) {
	a.mx.Lock()
	defer a.mx.Unlock()

	product, ok := a.productsByID[id]
	if !ok {
		return nil, persistence.ErrNotFound
	}

	return &model.Product{
		ID:    id,
		Name:  product.name,
		Price: product.price,
	}, nil
}

var _ persistence.CartRepository = (*Adapter)(nil)

type cart struct {
	userID    string
	positions []struct {
		ProductID string
		Quantity  int
		Price     int // in cents
	}
}

// CreateCart creates a cart for the given user with the given id and positions.
// Id must be unique. ErrConflict is returned otherwise.
func (a *Adapter) CreateCart(_ context.Context, userID, id string, positions []struct {
	ProductID string
	Quantity  int
	Price     int // in cents
}) error {
	a.mx.Lock()
	defer a.mx.Unlock()

	if _, ok := a.cartsByID[id]; ok {
		return persistence.ErrConflict
	}

	a.cartsByID[id] = &cart{
		userID:    userID,
		positions: positions,
	}
	return nil
}

// UpdateCartOfUser updates a cart of the given user with new positions. Any
// existing positions are replaced. ErrNotFound is returned if the cart does not
// exist. ErrNotOwnedByUser is returned if the cart exists but it's not owned by
// the given user.
func (a *Adapter) UpdateCartOfUser(ctx context.Context, userID, id string, positions []struct {
	ProductID string
	Quantity  int
	Price     int // in cents
}) error {
	a.mx.Lock()
	defer a.mx.Unlock()

	cart, ok := a.cartsByID[id]
	if !ok {
		return persistence.ErrNotFound
	}

	if cart.userID != userID {
		return persistence.ErrNotOwnedByUser
	}

	cart.positions = positions
	return nil
}

// FindAllUnlockedCartsOfUser returns all stored carts and their positions of
// the given user.
func (a *Adapter) FindAllUnlockedCartsOfUser(_ context.Context, userID string) ([]*model.Cart, error) {
	a.mx.Lock()
	defer a.mx.Unlock()

	result := make([]*model.Cart, 0)
	for id, cart := range a.cartsByID {
		if cart.userID != userID {
			continue
		}
		result = append(result, convertCartOut(id, cart))
	}
	return result, nil
}

// FindCartOfUser returns the cart of the given user with the given cart id.
// ErrNotFound is returned if there is no cart with the id. ErrNotOwnedByUser is
// returned if the cart exists but it's not owned by the given user.
func (a *Adapter) FindCartOfUser(_ context.Context, userID, id string) (*model.Cart, error) {
	a.mx.Lock()
	defer a.mx.Unlock()

	cart, ok := a.cartsByID[id]
	if !ok {
		return nil, persistence.ErrNotFound
	}

	if cart.userID != userID {
		return nil, persistence.ErrNotOwnedByUser
	}

	return convertCartOut(id, cart), nil
}

func convertCartOut(id string, cart *cart) *model.Cart {
	out := model.Cart{
		ID:        id,
		Positions: make([]model.Position, len(cart.positions)),
	}
	for i, position := range cart.positions {
		out.Positions[i] = model.Position{
			ProductID: position.ProductID,
			Quantity:  position.Quantity,
			Price:     position.Price, // TODO: don't save the price
		}
	}
	return &out
}
