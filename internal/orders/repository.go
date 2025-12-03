package orders

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository persists orders and items.
type Repository struct {
	pool *pgxpool.Pool
}

// NewRepository creates repository.
func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

var (
	errNilPool       = errors.New("orders repository: nil pool")
	errOrderNotFound = errors.New("order not found")
)

// CreateParams contains all fields to insert order + items.
type CreateParams struct {
	ClientRequestID string
	UserID          int64
	AddressID       *int64
	Type            string
	PaymentMethod   string
	Status          string
	RegionID        int64
	DeliveryPrice   float64
	ItemsTotal      float64
	TotalPrice      float64
	Comment         string
	CustomerName    string
	CustomerPhone   string
	Items           []OrderItem
}

// Create inserts order and items, returns complete record.
func (r *Repository) Create(ctx context.Context, params CreateParams) (*Order, error) {
	if r.pool == nil {
		return nil, errNilPool
	}

	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var addressID interface{}
	if params.AddressID != nil {
		addressID = *params.AddressID
	}

	row := tx.QueryRow(ctx, `
INSERT INTO orders (client_request_id, user_id, address_id, type, payment_method, status, region_id, delivery_price, items_total, total_price, comment, customer_name, customer_phone)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)
RETURNING id, client_request_id, user_id, COALESCE(address_id,0), type, payment_method, status, region_id, delivery_price, items_total, total_price, COALESCE(comment,''), customer_name, customer_phone, created_at, updated_at;
`,
		params.ClientRequestID,
		params.UserID,
		addressID,
		params.Type,
		params.PaymentMethod,
		params.Status,
		params.RegionID,
		params.DeliveryPrice,
		params.ItemsTotal,
		params.TotalPrice,
		params.Comment,
		params.CustomerName,
		params.CustomerPhone,
	)

	order := &Order{}
	if err := row.Scan(
		&order.ID,
		&order.ClientRequestID,
		&order.UserID,
		&order.AddressID,
		&order.Type,
		&order.PaymentMethod,
		&order.Status,
		&order.RegionID,
		&order.DeliveryPrice,
		&order.ItemsTotal,
		&order.TotalPrice,
		&order.Comment,
		&order.CustomerName,
		&order.CustomerPhone,
		&order.CreatedAt,
		&order.UpdatedAt,
	); err != nil {
		if isUniqueViolation(err) {
			existing, getErr := r.GetByClientRequestID(ctx, params.ClientRequestID, params.UserID)
			if getErr != nil {
				return nil, getErr
			}
			return existing, nil
		}
		return nil, err
	}

	batch := &pgx.Batch{}
	for _, item := range params.Items {
		batch.Queue(`
INSERT INTO order_items (order_id, product_id, product_name, qty, price, total)
VALUES ($1,$2,$3,$4,$5,$6)
RETURNING id;
`, order.ID, item.ProductID, item.ProductName, item.Qty, item.Price, item.Total)
	}

	br := tx.SendBatch(ctx, batch)
	order.Items = make([]OrderItem, len(params.Items))
	for i := range params.Items {
		var id int64
		if err := br.QueryRow().Scan(&id); err != nil {
			br.Close()
			return nil, err
		}
		order.Items[i] = params.Items[i]
		order.Items[i].ID = id
		order.Items[i].OrderID = order.ID
	}
	if err := br.Close(); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return order, nil
}

// GetByClientRequestID returns order + items for idempotency.
func (r *Repository) GetByClientRequestID(ctx context.Context, clientReqID string, userID int64) (*Order, error) {
	if r.pool == nil {
		return nil, errNilPool
	}

	order, err := r.fetchOrder(ctx, `client_request_id = $1 AND user_id = $2`, clientReqID, userID)
	if err != nil {
		return nil, err
	}
	return order, nil
}

// GetByID fetches order by id belonging to user.
func (r *Repository) GetByID(ctx context.Context, id, userID int64) (*Order, error) {
	if r.pool == nil {
		return nil, errNilPool
	}
	return r.fetchOrder(ctx, `id = $1 AND user_id = $2`, id, userID)
}

func (r *Repository) fetchOrder(ctx context.Context, where string, args ...interface{}) (*Order, error) {
	query := `
SELECT id, client_request_id, user_id, COALESCE(address_id,0), type, payment_method, status, region_id, delivery_price, items_total, total_price, COALESCE(comment,''), customer_name, customer_phone, created_at, updated_at
FROM orders
WHERE ` + where + `
LIMIT 1;
`

	row := r.pool.QueryRow(ctx, query, args...)
	var order Order
	if err := row.Scan(&order.ID, &order.ClientRequestID, &order.UserID, &order.AddressID, &order.Type, &order.PaymentMethod, &order.Status, &order.RegionID, &order.DeliveryPrice, &order.ItemsTotal, &order.TotalPrice, &order.Comment, &order.CustomerName, &order.CustomerPhone, &order.CreatedAt, &order.UpdatedAt); err != nil {
		return nil, err
	}

	items, err := r.fetchItems(ctx, order.ID)
	if err != nil {
		return nil, err
	}
	order.Items = items
	return &order, nil
}

// UpdateStatus updates order status and returns updated order.
func (r *Repository) UpdateStatus(ctx context.Context, orderID int64, status string) (*Order, error) {
	if r.pool == nil {
		return nil, errNilPool
	}
	const query = `
UPDATE orders SET status=$1, updated_at=NOW()
WHERE id=$2
RETURNING id, client_request_id, user_id, COALESCE(address_id,0), type, payment_method, status, region_id, delivery_price, items_total, total_price, COALESCE(comment,''), customer_name, customer_phone, created_at, updated_at;
`
	row := r.pool.QueryRow(ctx, query, status, orderID)
	var order Order
	if err := row.Scan(&order.ID, &order.ClientRequestID, &order.UserID, &order.AddressID, &order.Type, &order.PaymentMethod, &order.Status, &order.RegionID, &order.DeliveryPrice, &order.ItemsTotal, &order.TotalPrice, &order.Comment, &order.CustomerName, &order.CustomerPhone, &order.CreatedAt, &order.UpdatedAt); err != nil {
		return nil, err
	}
	items, err := r.fetchItems(ctx, order.ID)
	if err != nil {
		return nil, err
	}
	order.Items = items
	return &order, nil
}

// ListByUser returns latest orders for user.
func (r *Repository) ListByUser(ctx context.Context, userID int64, limit int) ([]Order, error) {
	if r.pool == nil {
		return nil, errNilPool
	}
	if limit <= 0 {
		limit = 50
	}

	const query = `
SELECT id, client_request_id, user_id, COALESCE(address_id,0), type, payment_method, status, region_id, delivery_price, items_total, total_price, COALESCE(comment,''), customer_name, customer_phone, created_at, updated_at
FROM orders
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2;
`

	rows, err := r.pool.Query(ctx, query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ordersList []Order
	for rows.Next() {
		var order Order
		if err := rows.Scan(&order.ID, &order.ClientRequestID, &order.UserID, &order.AddressID, &order.Type, &order.PaymentMethod, &order.Status, &order.RegionID, &order.DeliveryPrice, &order.ItemsTotal, &order.TotalPrice, &order.Comment, &order.CustomerName, &order.CustomerPhone, &order.CreatedAt, &order.UpdatedAt); err != nil {
			return nil, err
		}
		ordersList = append(ordersList, order)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	for i := range ordersList {
		items, err := r.fetchItems(ctx, ordersList[i].ID)
		if err != nil {
			return nil, err
		}
		ordersList[i].Items = items
	}

	return ordersList, nil
}

// ListAdmin returns orders with optional filters.
func (r *Repository) ListAdmin(ctx context.Context, params AdminListParams) ([]Order, error) {
	if r.pool == nil {
		return nil, errNilPool
	}
	limit := params.Limit
	if limit <= 0 {
		limit = 100
	}
	offset := params.Offset
	query := `
SELECT id, client_request_id, user_id, COALESCE(address_id,0), type, payment_method, status, region_id, delivery_price, items_total, total_price, COALESCE(comment,''), customer_name, customer_phone, created_at, updated_at
FROM orders
WHERE 1=1`
	args := []interface{}{}
	idx := 1
	if params.Status != "" {
		query += fmt.Sprintf(" AND status = $%d", idx)
		args = append(args, params.Status)
		idx++
	}
	if params.From != nil {
		query += fmt.Sprintf(" AND created_at >= $%d", idx)
		args = append(args, *params.From)
		idx++
	}
	if params.To != nil {
		query += fmt.Sprintf(" AND created_at <= $%d", idx)
		args = append(args, *params.To)
		idx++
	}
	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", idx, idx+1)
	args = append(args, limit, offset)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ordersList []Order
	for rows.Next() {
		var order Order
		if err := rows.Scan(&order.ID, &order.ClientRequestID, &order.UserID, &order.AddressID, &order.Type, &order.PaymentMethod, &order.Status, &order.RegionID, &order.DeliveryPrice, &order.ItemsTotal, &order.TotalPrice, &order.Comment, &order.CustomerName, &order.CustomerPhone, &order.CreatedAt, &order.UpdatedAt); err != nil {
			return nil, err
		}
		ordersList = append(ordersList, order)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	for i := range ordersList {
		items, err := r.fetchItems(ctx, ordersList[i].ID)
		if err != nil {
			return nil, err
		}
		ordersList[i].Items = items
	}

	return ordersList, nil
}

func (r *Repository) fetchItems(ctx context.Context, orderID int64) ([]OrderItem, error) {
	const query = `
SELECT id, order_id, product_id, product_name, qty, price, total
FROM order_items
WHERE order_id = $1;
`

	rows, err := r.pool.Query(ctx, query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []OrderItem
	for rows.Next() {
		var oi OrderItem
		if err := rows.Scan(&oi.ID, &oi.OrderID, &oi.ProductID, &oi.ProductName, &oi.Qty, &oi.Price, &oi.Total); err != nil {
			return nil, err
		}
		items = append(items, oi)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}
