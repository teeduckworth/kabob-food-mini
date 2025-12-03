package menu

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository loads menu data.
type Repository struct {
	pool *pgxpool.Pool
}

// NewRepository returns a new menu repository.
func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

var errNilPool = errors.New("menu repository: nil pool")

// InsertCategory creates a category.
func (r *Repository) InsertCategory(ctx context.Context, cat Category) (*Category, error) {
	if r.pool == nil {
		return nil, errNilPool
	}
	const query = `
INSERT INTO categories (name, emoji, sort_order, is_active)
VALUES ($1,$2,$3,$4)
RETURNING id, name, emoji, sort_order, is_active;
`
	row := r.pool.QueryRow(ctx, query, cat.Name, cat.Emoji, cat.SortOrder, cat.IsActive)
	var res Category
	if err := row.Scan(&res.ID, &res.Name, &res.Emoji, &res.SortOrder, &res.IsActive); err != nil {
		return nil, err
	}
	return &res, nil
}

// UpdateCategory updates category values.
func (r *Repository) UpdateCategory(ctx context.Context, cat Category) (*Category, error) {
	if r.pool == nil {
		return nil, errNilPool
	}
	const query = `
UPDATE categories
SET name=$1, emoji=$2, sort_order=$3, is_active=$4
WHERE id=$5
RETURNING id, name, emoji, sort_order, is_active;
`
	row := r.pool.QueryRow(ctx, query, cat.Name, cat.Emoji, cat.SortOrder, cat.IsActive, cat.ID)
	var res Category
	if err := row.Scan(&res.ID, &res.Name, &res.Emoji, &res.SortOrder, &res.IsActive); err != nil {
		return nil, err
	}
	return &res, nil
}

// DeleteCategory removes category.
func (r *Repository) DeleteCategory(ctx context.Context, id int64) error {
	if r.pool == nil {
		return errNilPool
	}
	const query = `DELETE FROM categories WHERE id=$1;`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}

// InsertProduct creates a product.
func (r *Repository) InsertProduct(ctx context.Context, product Product) (*Product, error) {
	if r.pool == nil {
		return nil, errNilPool
	}
	const query = `
INSERT INTO products (category_id, name, description, price, old_price, image_url, is_active, sort_order)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
RETURNING id, category_id, name, description, price, COALESCE(old_price,0), image_url, is_active, sort_order;
`
	row := r.pool.QueryRow(ctx, query, product.CategoryID, product.Name, product.Description, product.Price, product.OldPrice, product.ImageURL, product.IsActive, product.SortOrder)
	var res Product
	if err := row.Scan(&res.ID, &res.CategoryID, &res.Name, &res.Description, &res.Price, &res.OldPrice, &res.ImageURL, &res.IsActive, &res.SortOrder); err != nil {
		return nil, err
	}
	return &res, nil
}

// UpdateProduct updates values.
func (r *Repository) UpdateProduct(ctx context.Context, product Product) (*Product, error) {
	if r.pool == nil {
		return nil, errNilPool
	}
	const query = `
UPDATE products
SET category_id=$1,
    name=$2,
    description=$3,
    price=$4,
    old_price=$5,
    image_url=$6,
    is_active=$7,
    sort_order=$8
WHERE id=$9
RETURNING id, category_id, name, description, price, COALESCE(old_price,0), image_url, is_active, sort_order;
`
	row := r.pool.QueryRow(ctx, query, product.CategoryID, product.Name, product.Description, product.Price, product.OldPrice, product.ImageURL, product.IsActive, product.SortOrder, product.ID)
	var res Product
	if err := row.Scan(&res.ID, &res.CategoryID, &res.Name, &res.Description, &res.Price, &res.OldPrice, &res.ImageURL, &res.IsActive, &res.SortOrder); err != nil {
		return nil, err
	}
	return &res, nil
}

// DeleteProduct removes product by id.
func (r *Repository) DeleteProduct(ctx context.Context, id int64) error {
	if r.pool == nil {
		return errNilPool
	}
	const query = `DELETE FROM products WHERE id=$1;`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}

// GetActiveCategories returns all active categories ordered by sort_order.
func (r *Repository) GetActiveCategories(ctx context.Context) ([]Category, error) {
	if r.pool == nil {
		return nil, errNilPool
	}

	const query = `
SELECT id, name, emoji, sort_order, is_active
FROM categories
WHERE is_active = TRUE
ORDER BY sort_order ASC, id;
`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var cat Category
		if err := rows.Scan(&cat.ID, &cat.Name, &cat.Emoji, &cat.SortOrder, &cat.IsActive); err != nil {
			return nil, err
		}
		categories = append(categories, cat)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

// GetActiveProducts returns active products ordered by category -> sort_order.
func (r *Repository) GetActiveProducts(ctx context.Context) ([]Product, error) {
	if r.pool == nil {
		return nil, errNilPool
	}

	const query = `
SELECT id, category_id, name, description, price, COALESCE(old_price, 0), image_url, is_active, sort_order
FROM products
WHERE is_active = TRUE
ORDER BY category_id, sort_order ASC, id;
`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.CategoryID, &p.Name, &p.Description, &p.Price, &p.OldPrice, &p.ImageURL, &p.IsActive, &p.SortOrder); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}
