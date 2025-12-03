package addresses

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository works with address records.
type Repository struct {
	pool *pgxpool.Pool
}

// NewRepository creates repository.
func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

var errNilPool = errors.New("addresses repository: nil pool")

// ListByUser returns addresses belonging to user sorted by created_at desc.
func (r *Repository) ListByUser(ctx context.Context, userID int64) ([]Address, error) {
	if r.pool == nil {
		return nil, errNilPool
	}

	const query = `
SELECT id, user_id, region_id, street, house, COALESCE(entrance, ''), COALESCE(flat, ''), COALESCE(comment, ''), is_default, created_at
FROM addresses
WHERE user_id = $1
ORDER BY is_default DESC, created_at DESC;
`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []Address
	for rows.Next() {
		var addr Address
		if err := rows.Scan(&addr.ID, &addr.UserID, &addr.RegionID, &addr.Street, &addr.House, &addr.Entrance, &addr.Flat, &addr.Comment, &addr.IsDefault, &addr.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, addr)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
}

// GetByIDAndUser returns address ensuring ownership.
func (r *Repository) GetByIDAndUser(ctx context.Context, id, userID int64) (*Address, error) {
	if r.pool == nil {
		return nil, errNilPool
	}

	const query = `
SELECT id, user_id, region_id, street, house, COALESCE(entrance, ''), COALESCE(flat, ''), COALESCE(comment, ''), is_default, created_at
FROM addresses
WHERE id = $1 AND user_id = $2;
`

	row := r.pool.QueryRow(ctx, query, id, userID)

	var addr Address
	if err := row.Scan(
		&addr.ID,
		&addr.UserID,
		&addr.RegionID,
		&addr.Street,
		&addr.House,
		&addr.Entrance,
		&addr.Flat,
		&addr.Comment,
		&addr.IsDefault,
		&addr.CreatedAt,
	); err != nil {
		return nil, err
	}

	return &addr, nil
}

// Insert creates new address.
func (r *Repository) Insert(ctx context.Context, input CreateInput) (*Address, error) {
	if r.pool == nil {
		return nil, errNilPool
	}

	const query = `
INSERT INTO addresses (user_id, region_id, street, house, entrance, flat, comment, is_default)
VALUES ($1, $2, $3, $4, NULLIF($5, ''), NULLIF($6, ''), NULLIF($7, ''), $8)
RETURNING id, user_id, region_id, street, house, COALESCE(entrance, ''), COALESCE(flat, ''), COALESCE(comment, ''), is_default, created_at;
`

	row := r.pool.QueryRow(ctx, query,
		input.UserID,
		input.RegionID,
		input.Street,
		input.House,
		input.Entrance,
		input.Flat,
		input.Comment,
		input.IsDefault,
	)

	var addr Address
	if err := row.Scan(
		&addr.ID,
		&addr.UserID,
		&addr.RegionID,
		&addr.Street,
		&addr.House,
		&addr.Entrance,
		&addr.Flat,
		&addr.Comment,
		&addr.IsDefault,
		&addr.CreatedAt,
	); err != nil {
		return nil, err
	}

	return &addr, nil
}

// Update modifies user address.
func (r *Repository) Update(ctx context.Context, input UpdateInput) (*Address, error) {
	if r.pool == nil {
		return nil, errNilPool
	}

	const query = `
UPDATE addresses
SET region_id = $1,
    street = $2,
    house = $3,
    entrance = NULLIF($4, ''),
    flat = NULLIF($5, ''),
    comment = NULLIF($6, ''),
    is_default = $7
WHERE id = $8 AND user_id = $9
RETURNING id, user_id, region_id, street, house, COALESCE(entrance, ''), COALESCE(flat, ''), COALESCE(comment, ''), is_default, created_at;
`

	row := r.pool.QueryRow(ctx, query,
		input.RegionID,
		input.Street,
		input.House,
		input.Entrance,
		input.Flat,
		input.Comment,
		input.IsDefault,
		input.ID,
		input.UserID,
	)

	var addr Address
	if err := row.Scan(
		&addr.ID,
		&addr.UserID,
		&addr.RegionID,
		&addr.Street,
		&addr.House,
		&addr.Entrance,
		&addr.Flat,
		&addr.Comment,
		&addr.IsDefault,
		&addr.CreatedAt,
	); err != nil {
		return nil, err
	}

	return &addr, nil
}

// Delete removes an address belonging to user.
func (r *Repository) Delete(ctx context.Context, id, userID int64) error {
	if r.pool == nil {
		return errNilPool
	}

	const query = `DELETE FROM addresses WHERE id = $1 AND user_id = $2;`
	cmd, err := r.pool.Exec(ctx, query, id, userID)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("address not found")
	}
	return nil
}

// ClearDefault unsets default flag for all user addresses.
func (r *Repository) ClearDefault(ctx context.Context, userID int64) error {
	if r.pool == nil {
		return errNilPool
	}
	const query = `UPDATE addresses SET is_default = FALSE WHERE user_id = $1;`
	_, err := r.pool.Exec(ctx, query, userID)
	return err
}
