package regions

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository loads region data from PostgreSQL.
type Repository struct {
	pool *pgxpool.Pool
}

// NewRepository constructs region repo.
func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

var errNilPool = errors.New("regions repository: nil pool")

// Insert adds a region.
func (r *Repository) Insert(ctx context.Context, region Region) (*Region, error) {
	if r.pool == nil {
		return nil, errNilPool
	}
	const query = `
INSERT INTO regions (name, delivery_price, is_active)
VALUES ($1,$2,$3)
RETURNING id, name, delivery_price, is_active;
`
	row := r.pool.QueryRow(ctx, query, region.Name, region.DeliveryPrice, region.IsActive)
	var res Region
	if err := row.Scan(&res.ID, &res.Name, &res.DeliveryPrice, &res.IsActive); err != nil {
		return nil, err
	}
	return &res, nil
}

// Update modifies region.
func (r *Repository) Update(ctx context.Context, region Region) (*Region, error) {
	if r.pool == nil {
		return nil, errNilPool
	}
	const query = `
UPDATE regions SET name=$1, delivery_price=$2, is_active=$3
WHERE id=$4
RETURNING id, name, delivery_price, is_active;
`
	row := r.pool.QueryRow(ctx, query, region.Name, region.DeliveryPrice, region.IsActive, region.ID)
	var res Region
	if err := row.Scan(&res.ID, &res.Name, &res.DeliveryPrice, &res.IsActive); err != nil {
		return nil, err
	}
	return &res, nil
}

// Delete removes region.
func (r *Repository) Delete(ctx context.Context, id int64) error {
	if r.pool == nil {
		return errNilPool
	}
	const query = `DELETE FROM regions WHERE id = $1;`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}

// GetActiveRegions returns currently active regions sorted by name.
func (r *Repository) GetActiveRegions(ctx context.Context) ([]Region, error) {
	if r.pool == nil {
		return nil, errNilPool
	}

	const query = `
SELECT id, name, delivery_price, is_active
FROM regions
WHERE is_active = TRUE
ORDER BY id;
`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []Region
	for rows.Next() {
		var region Region
		if err := rows.Scan(&region.ID, &region.Name, &region.DeliveryPrice, &region.IsActive); err != nil {
			return nil, err
		}
		result = append(result, region)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

// GetByID returns region by id.
func (r *Repository) GetByID(ctx context.Context, id int64) (*Region, error) {
	if r.pool == nil {
		return nil, errNilPool
	}

	const query = `
SELECT id, name, delivery_price, is_active
FROM regions
WHERE id = $1;
`

	row := r.pool.QueryRow(ctx, query, id)
	var region Region
	if err := row.Scan(&region.ID, &region.Name, &region.DeliveryPrice, &region.IsActive); err != nil {
		return nil, err
	}
	return &region, nil
}
