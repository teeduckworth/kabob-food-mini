package menu

// Category represents menu category metadata.
type Category struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Emoji     string `json:"emoji"`
	SortOrder int    `json:"sort_order"`
	IsActive  bool   `json:"is_active"`
}

// Product represents an item that belongs to a category.
type Product struct {
	ID          int64   `json:"id"`
	CategoryID  int64   `json:"category_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	OldPrice    float64 `json:"old_price"`
	ImageURL    string  `json:"image_url"`
	IsActive    bool    `json:"is_active"`
	SortOrder   int     `json:"sort_order"`
}

// MenuCategory wraps category info with its products.
type MenuCategory struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Emoji     string    `json:"emoji"`
	SortOrder int       `json:"sort_order"`
	Products  []Product `json:"products"`
}

// MenuResponse holds categories with products for /menu.
type MenuResponse struct {
	Categories []MenuCategory `json:"categories"`
}
