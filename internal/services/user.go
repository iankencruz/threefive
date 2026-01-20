package services

import (
	"github.com/iankencruz/threefive/database/generated"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserService struct {
	db      *pgxpool.Pool
	queries *generated.Queries
}

// NewUserService creates a new user service
func NewUserService(db *pgxpool.Pool, queries *generated.Queries) *UserService {
	return &UserService{
		db:      db,
		queries: queries,
	}
}
