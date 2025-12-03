package users

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository provides user persistence helpers.
type Repository struct {
	pool *pgxpool.Pool
}

// NewRepository returns a Postgres-backed Repository.
func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

var errNilPool = errors.New("nil pgx pool")

// UpsertTelegramUser inserts or updates a user based on Telegram profile data.
func (r *Repository) UpsertTelegramUser(ctx context.Context, input UpsertTelegramUserInput) (*User, error) {
	if r.pool == nil {
		return nil, errNilPool
	}

	query := `
INSERT INTO users (telegram_id, first_name, last_name, username, phone, language)
VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT (telegram_id) DO UPDATE SET
    first_name = EXCLUDED.first_name,
    last_name = EXCLUDED.last_name,
    username = EXCLUDED.username,
    phone = COALESCE(NULLIF(EXCLUDED.phone, ''), users.phone),
    language = EXCLUDED.language
RETURNING id, telegram_id, first_name, last_name, username, phone, language, created_at;
`

	row := r.pool.QueryRow(ctx, query,
		input.TelegramID,
		input.FirstName,
		input.LastName,
		input.Username,
		input.Phone,
		input.Language,
	)

	var user User
	if err := row.Scan(
		&user.ID,
		&user.TelegramID,
		&user.FirstName,
		&user.LastName,
		&user.Username,
		&user.Phone,
		&user.Language,
		&user.CreatedAt,
	); err != nil {
		return nil, err
	}

	return &user, nil
}

// GetByID fetches a user by id.
func (r *Repository) GetByID(ctx context.Context, id int64) (*User, error) {
	if r.pool == nil {
		return nil, errNilPool
	}

	const query = `
SELECT id, telegram_id, first_name, last_name, username, phone, language, created_at
FROM users
WHERE id = $1;
`

	row := r.pool.QueryRow(ctx, query, id)
	var user User
	if err := row.Scan(
		&user.ID,
		&user.TelegramID,
		&user.FirstName,
		&user.LastName,
		&user.Username,
		&user.Phone,
		&user.Language,
		&user.CreatedAt,
	); err != nil {
		return nil, err
	}

	return &user, nil
}
