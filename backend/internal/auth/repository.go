package auth

import (
	"context"

	"github.com/iankencruz/threefive/backend/internal/generated"
)

type Repository interface {
	CreateUser(ctx context.Context, arg generated.CreateUserParams) (*generated.User, error)
	GetUserByEmail(ctx context.Context, email string) (*generated.User, error)
	GetUserByID(ctx context.Context, id int32) (*generated.User, error)
	DeleteUserByID(ctx context.Context, id int32) error
}

type AuthRepository struct {
	q *generated.Queries
}

func NewAuthRepository(q *generated.Queries) Repository {
	return &AuthRepository{
		q: q,
	}
}

func (r *AuthRepository) CreateUser(ctx context.Context, arg generated.CreateUserParams) (*generated.User, error) {
	user, err := r.q.CreateUser(ctx, arg)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *AuthRepository) GetUserByEmail(ctx context.Context, email string) (*generated.User, error) {
	user, err := r.q.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *AuthRepository) GetUserByID(ctx context.Context, id int32) (*generated.User, error) {
	user, err := r.q.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *AuthRepository) DeleteUserByID(ctx context.Context, id int32) error {
	return r.q.DeleteUser(ctx, id)
}
