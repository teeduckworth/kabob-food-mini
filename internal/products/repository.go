package products

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository fetches product data.
type Repository struct {
	pool *pgxpool.Pool
}

// NewRepository creates repository.
func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

var errNilPool = errors.New("products repository: nil pool")

// GetActiveByIDs returns active products with ids.
func (r *Repository) GetActiveByIDs(ctx context.Context, ids []int64) (map[int64]Product, error) {
	if r.pool == nil {
		return nil, errNilPool
	}
	if len(ids) == 0 {
		return map[int64]Product{}, nil
	}

	const query = `
SELECT id, category_id, name, description, price, is_active
FROM products
WHERE id = ANY($1)
`

	rows, err := r.pool.Query(ctx, query, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[int64]Product)
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.CategoryID, &p.Name, &p.Description, &p.Price, &p.IsActive); err != nil {
			return nil, err
		}
		result[p.ID] = p
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
