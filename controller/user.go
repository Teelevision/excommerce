package controller

import (
	"context"
	"errors"
	"fmt"

	"github.com/Teelevision/excommerce/model"
	"github.com/Teelevision/excommerce/persistence"
	"github.com/google/uuid"
)

// User is the controller that creates a user.
type User struct {
	UserRepository persistence.UserRepository
}

// Create creates the user. The name is expected to be 1 to 64 runes long, and
// the password 8 to 64. ErrConflict is returned if the name is already taken.
// On success the user is returned.
func (c *User) Create(ctx context.Context, name, password string) (*model.User, error) {
	// create id
	uuid, err := uuid.NewRandom()
	if err != nil {
		panic(err)
	}
	id := uuid.String()

	// store user
	err = c.UserRepository.CreateUser(ctx, id, name, password)
	switch {
	case errors.Is(err, persistence.ErrConflict):
		return nil, fmt.Errorf("%w: %s", ErrConflict, err)
	case errors.Is(err, context.Canceled), errors.Is(err, context.DeadlineExceeded):
		return nil, err
	case err == nil:
		return &model.User{
			ID:   id,
			Name: name,
		}, nil
	default:
		panic(err)
	}
}

// GetByNameAndPassword gets the user by name and password. ErrNotFound is
// returned if there is no user with the name. On success the user is returned.
func (c *User) GetByNameAndPassword(ctx context.Context, name, password string) (*model.User, error) {
	user, err := c.UserRepository.FindUserByNameAndPassword(ctx, name, password)
	switch {
	case errors.Is(err, persistence.ErrNotFound):
		return nil, fmt.Errorf("%w: %s", ErrNotFound, err)
	case errors.Is(err, context.Canceled), errors.Is(err, context.DeadlineExceeded):
		return nil, err
	case err == nil:
		return user, nil
	default:
		panic(err)
	}
}
