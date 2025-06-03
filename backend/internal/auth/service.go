package auth

import (
	"context"
)

type Service interface {
	Register(ctx context.Context, firstName, lastName, email, password string) (*User, error)
	Login(ctx context.Context, email, password string) (*User, error)
	GetUserByID(ctx context.Context, id int32) (*User, error)
}

type AuthService struct {
	Repo Repository
}

func NewService(repo Repository) *AuthService {
	return &AuthService{Repo: repo}
}

func (s *AuthService) Register(ctx context.Context, firstName, lastName, email, password string) (*User, error) {
	// TODO: implement registration logic
	return nil, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*User, error) {
	// TODO: implement login logic
	return nil, nil
}

func (s *AuthService) GetUserByID(ctx context.Context, id int32) (*User, error) {
	return s.Repo.GetUserByID(ctx, id)
}
