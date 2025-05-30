package users

import (
	"context"

	"github.com/iankencruz/threefive/internal/models/users/sqlc"
)

type User = sqlc.User
type CreateUserParams = sqlc.CreateUserParams

type Repository interface {
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	ListUsers(ctx context.Context) ([]User, error)
}

type repository struct {
	q *sqlc.Queries
}

func NewRepository(db sqlc.DBTX) Repository {
	return &repository{
		q: sqlc.New(db),
	}
}

func (r *repository) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	return r.q.CreateUser(ctx, arg)
}

func (r *repository) GetUserByEmail(ctx context.Context, email string) (User, error) {
	return r.q.GetUserByEmail(ctx, email)
}

func (r *repository) ListUsers(ctx context.Context) ([]User, error) {
	return r.q.ListUsers(ctx)
}
