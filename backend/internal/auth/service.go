package auth

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Register(ctx context.Context, firstName, lastName, email, password string) (*User, error)
	Login(ctx context.Context, email, password string) (*User, error)
	GetUserByID(ctx context.Context, id int32) (*User, error)
}

type RegisterParams struct {
	FirstName string
	LastName  string
	Email     string
	Password  string
}

type AuthService struct {
	Repo Repository
}

func NewService(repo Repository) *AuthService {
	return &AuthService{Repo: repo}
}

func (s *AuthService) Register(ctx context.Context, firstName, lastName, email, password string) (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("could not hash password")
	}

	args := CreateUserParams{
		FirstName:    firstName,
		LastName:     lastName,
		Email:        email,
		PasswordHash: string(hashedPassword),
	}

	return s.Repo.CreateUser(ctx, args)
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*User, error) {
	user, err := s.Repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if !CheckPasswordHash(password, user.PasswordHash) {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

func (s *AuthService) GetUserByID(ctx context.Context, id int32) (*User, error) {
	return s.Repo.GetUserByID(ctx, id)
}

// === Utilities ===//
// HashPassword hashes a plain-text password for storage.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash compares a hashed password with a plain-text candidate.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
