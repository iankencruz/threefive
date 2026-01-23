// templates/lib/context.go
package lib

import (
	"context"

	"github.com/iankencruz/threefive/database/generated"
)

type contextKey string

const userContextKey contextKey = "user"

// WithUser adds the user to the context
func WithUser(ctx context.Context, user *generated.User) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}

// GetUser retrieves the user from the context
func GetUser(ctx context.Context) *generated.User {
	user, ok := ctx.Value(userContextKey).(*generated.User)
	if !ok {
		return nil
	}
	return user
}
