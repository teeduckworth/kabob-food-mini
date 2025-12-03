package products

// Product represents a product with price info used for orders.
type Product struct {
	ID          int64   `json:"id"`
	CategoryID  int64   `json:"category_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	IsActive    bool    `json:"is_active"`
}
