package addresses

import "time"

// Address describes delivery address data tied to a user.
type Address struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	RegionID  int64     `json:"region_id"`
	Street    string    `json:"street"`
	House     string    `json:"house"`
	Entrance  string    `json:"entrance"`
	Flat      string    `json:"flat"`
	Comment   string    `json:"comment"`
	IsDefault bool      `json:"is_default"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateInput defines request payload for new address.
type CreateInput struct {
	UserID    int64
	RegionID  int64
	Street    string
	House     string
	Entrance  string
	Flat      string
	Comment   string
	IsDefault bool
}

// UpdateInput defines payload for updating address fields.
type UpdateInput struct {
	ID        int64
	UserID    int64
	RegionID  int64
	Street    string
	House     string
	Entrance  string
	Flat      string
	Comment   string
	IsDefault bool
}
