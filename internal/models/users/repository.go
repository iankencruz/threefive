package models

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
	GetUserByID(ctx context.Context, id int32) (User, error)
}

type UserRepo struct {
	q *sqlc.Queries
}

func NewUser(db sqlc.DBTX) Repository {
	return &UserRepo{
		q: sqlc.New(db),
	}
}

func (r *UserRepo) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	return r.q.CreateUser(ctx, arg)
}

func (r *UserRepo) GetUserByEmail(ctx context.Context, email string) (User, error) {
	return r.q.GetUserByEmail(ctx, email)
}

func (r *UserRepo) ListUsers(ctx context.Context) ([]User, error) {
	return r.q.ListUsers(ctx)
}

func (r *UserRepo) GetUserByID(ctx context.Context, id int32) (User, error) {
	return r.q.GetUserByID(ctx, id)
}
