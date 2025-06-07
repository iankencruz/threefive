package auth

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	CreateUser(ctx context.Context, arg CreateUserParams) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByID(ctx context.Context, id int32) (*User, error)
}

type PgxRepository struct {
	q *Queries
}

func NewPgxRepository(pool *pgxpool.Pool) *PgxRepository {
	return &PgxRepository{
		q: New(pool),
	}
}

func (r *PgxRepository) CreateUser(ctx context.Context, arg CreateUserParams) (*User, error) {
	if strings.TrimSpace(arg.Email) == "" {
		return nil, errors.New("email is required")
	}
	user, err := r.q.CreateUser(ctx, arg)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *PgxRepository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	user, err := r.q.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *PgxRepository) GetUserByID(ctx context.Context, id int32) (*User, error) {
	user, err := r.q.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *PgxRepository) DeleteUserByID(ctx context.Context, id int32) error {
	err := r.q.DeleteUser(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
