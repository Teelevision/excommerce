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

	usersByID   map[string]*user
	usersByName map[string]*user

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
		usersByID:   make(map[string]*user),
		usersByName: make(map[string]*user),
	}
	for _, option := range options {
		option(&a)
	}
	return &a
}

// The repositories that the in-memory adapter implements.
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
