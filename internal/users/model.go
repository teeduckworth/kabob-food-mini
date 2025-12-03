package users

import "time"

// User represents an application user persisted in PostgreSQL.
type User struct {
	ID         int64     `json:"id"`
	TelegramID int64     `json:"telegram_id"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	Username   string    `json:"username"`
	Phone      string    `json:"phone"`
	Language   string    `json:"language"`
	CreatedAt  time.Time `json:"created_at"`
}

// UpsertTelegramUserInput carries Telegram profile info to persist/update user.
type UpsertTelegramUserInput struct {
	TelegramID int64
	FirstName  string
	LastName   string
	Username   string
	Phone      string
	Language   string
}
