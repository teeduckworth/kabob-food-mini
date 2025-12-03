package regions

// Region describes delivery region metadata.
type Region struct {
	ID            int64   `json:"id"`
	Name          string  `json:"name"`
	DeliveryPrice float64 `json:"delivery_price"`
	IsActive      bool    `json:"is_active"`
}
