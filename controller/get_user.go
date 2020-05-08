package controller

import (
	"context"
	"errors"
	"fmt"

	"github.com/Teelevision/excommerce/model"
	"github.com/Teelevision/excommerce/persistence"
)

// GetUser is the controller that gets a user.
type GetUser struct {
	UserRepository persistence.UserRepository
}

// ByNameAndPassword gets the user by name and password. ErrNotFound is returned
// if there is no user with the name. On success the user is returned.
func (c *GetUser) ByNameAndPassword(ctx context.Context, name, password string) (*model.User, error) {
	user, err := c.UserRepository.FindUserByNameAndPassword(ctx, name, password)
	switch {
	case errors.Is(err, persistence.ErrNotFound):
		return nil, fmt.Errorf("%w: %s", ErrNotFound, err)
	case err == nil:
		return user, nil
	default:
		panic(err)
	}
}

