package authentication

import (
	"context"
	"errors"
	"net/http"

	"github.com/Teelevision/excommerce/model"
	"github.com/Teelevision/excommerce/persistence"
)

// Authenticator authenticates users. If used as a middleware it requires that
// the request is authenticated.
type Authenticator struct {
	UserRepository persistence.UserRepository
}

type userCtxKey struct{}

// AuthenticatedUser returns the authenticated user in the context or nil if no
// user is authenticated. Modifying the returned user does not modify the user
// in the context.
func AuthenticatedUser(ctx context.Context) *model.User {
	v := ctx.Value(userCtxKey{})
	if v == nil {
		return nil
	}
	user, ok := v.(model.User)
	if !ok {
		return nil
	}
	return &user
}

// Middleware returns a handler that authenticates the user making the request
// and add that user to the context. Using this middleware enables the usage of
// AuthenticatedUser to retrieve the user that made the request.
func (a *Authenticator) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, password, ok := r.BasicAuth()
		if !ok {
			w.WriteHeader(http.StatusUnauthorized) // 401
			return
		}

		ctx := r.Context()

		user, err := a.UserRepository.FindUserByIDAndPassword(ctx, id, password)
		switch {
		case errors.Is(err, persistence.ErrNotFound):
			w.WriteHeader(http.StatusUnauthorized) // 401
		case err == nil:
			ctx = context.WithValue(ctx, userCtxKey{}, *user)
			next.ServeHTTP(w, r.WithContext(ctx))
		default:
			panic(err)
		}
	})
}
