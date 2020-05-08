package controller

import (
	"context"
	"errors"
	"fmt"

	"github.com/Teelevision/excommerce/model"
	"github.com/Teelevision/excommerce/persistence"
	"github.com/google/uuid"
)

// CreateUser is the controller that creates a user.
type CreateUser struct {
	UserRepository persistence.UserRepository
}

// Do creates the user. The name is expected to be 1 to 64 runes long, and the
// password 8 to 64. ErrConflict is returned if the name is already taken. On
// success the user is returned.
func (c *CreateUser) Do(ctx context.Context, name, password string) (*model.User, error) {
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
	case err == nil:
		return &model.User{
			ID:   id,
			Name: name,
		}, nil
	default:
		panic(err)
	}
}
