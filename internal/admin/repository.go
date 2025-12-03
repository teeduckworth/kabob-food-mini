package admin

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository handles admin user persistence.
type Repository struct {
	pool *pgxpool.Pool
}

// NewRepository builds repository.
func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

var errNilPool = errors.New("admin repository: nil pool")

// EnsureUser ensures username exists else creates with provided hash.
func (r *Repository) EnsureUser(ctx context.Context, username, passwordHash string) error {
	if r.pool == nil {
		return errNilPool
	}
	const query = `
INSERT INTO admin_users (username, password_hash)
VALUES ($1, $2)
ON CONFLICT (username) DO NOTHING;
`
	_, err := r.pool.Exec(ctx, query, username, passwordHash)
	return err
}

// GetByUsername returns admin user by username.
func (r *Repository) GetByUsername(ctx context.Context, username string) (*User, error) {
	if r.pool == nil {
		return nil, errNilPool
	}

	const query = `SELECT id, username, password_hash, created_at FROM admin_users WHERE username = $1;`

	row := r.pool.QueryRow(ctx, query, username)
	var user User
	if err := row.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.CreatedAt); err != nil {
		return nil, err
	}
	return &user, nil
}
