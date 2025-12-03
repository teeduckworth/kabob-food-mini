package admin

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// AuthService handles admin authentication.
type AuthService struct {
	repo      *Repository
	jwtSecret []byte
}

// AuthConfig contains dependencies.
type AuthConfig struct {
	Repo      *Repository
	JWTSecret string
}

// NewAuthService builds service.
func NewAuthService(cfg AuthConfig) (*AuthService, error) {
	if cfg.Repo == nil {
		return nil, errors.New("admin repo is required")
	}
	if cfg.JWTSecret == "" {
		return nil, errors.New("jwt secret required")
	}
	return &AuthService{repo: cfg.Repo, jwtSecret: []byte(cfg.JWTSecret)}, nil
}

// Login authenticates admin credentials.
func (s *AuthService) Login(ctx context.Context, username, password string) (*User, error) {
	user, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}
	return user, nil
}

// EnsureDefaultAdmin creates an admin with provided credentials if missing.
func (s *AuthService) EnsureDefaultAdmin(ctx context.Context, username, password string) error {
	if username == "" || password == "" {
		return errors.New("default admin credentials not provided")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return s.repo.EnsureUser(ctx, username, string(hash))
}
