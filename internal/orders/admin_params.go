package orders

import "time"

// AdminListParams defines filters for admin order listing.
type AdminListParams struct {
	Status string
	From   *time.Time
	To     *time.Time
	Limit  int
	Offset int
}
