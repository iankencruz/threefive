package application

import (
	"github.com/iankencruz/threefive/internal/models/users"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Application struct {
	Users users.Repository
	DB    *pgxpool.Pool
}

func NewApp(db *pgxpool.Pool) *Application {

	return &Application{
		DB:    db,
		Users: users.NewRepository(db),
	}
}
