package orders

import "time"

// Order represents persisted order with items.
type Order struct {
	ID              int64       `json:"id"`
	ClientRequestID string      `json:"client_request_id"`
	UserID          int64       `json:"user_id"`
	AddressID       int64       `json:"address_id"`
	Type            string      `json:"type"`
	PaymentMethod   string      `json:"payment_method"`
	Status          string      `json:"status"`
	RegionID        int64       `json:"region_id"`
	DeliveryPrice   float64     `json:"delivery_price"`
	ItemsTotal      float64     `json:"items_total"`
	TotalPrice      float64     `json:"total_price"`
	Comment         string      `json:"comment"`
	CustomerName    string      `json:"customer_name"`
	CustomerPhone   string      `json:"customer_phone"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
	Items           []OrderItem `json:"items"`
}

// OrderItem represents product snapshot inside order.
type OrderItem struct {
	ID          int64   `json:"id"`
	OrderID     int64   `json:"order_id"`
	ProductID   int64   `json:"product_id"`
	ProductName string  `json:"product_name"`
	Qty         int32   `json:"qty"`
	Price       float64 `json:"price"`
	Total       float64 `json:"total"`
}

// ItemInput from API request.
type ItemInput struct {
	ProductID int64 `json:"product_id"`
	Qty       int32 `json:"qty"`
}

// CreateOrderInput from user request.
type CreateOrderInput struct {
	ClientRequestID string      `json:"client_request_id"`
	Type            string      `json:"type"`
	RegionID        int64       `json:"region_id"`
	AddressID       int64       `json:"address_id"`
	PaymentMethod   string      `json:"payment_method"`
	CustomerName    string      `json:"customer_name"`
	CustomerPhone   string      `json:"customer_phone"`
	Comment         string      `json:"comment"`
	Items           []ItemInput `json:"items"`
}

// UpdateStatusInput used by admin to change order status.
type UpdateStatusInput struct {
	Status string `json:"status" binding:"required"`
}
