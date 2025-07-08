package models

import "time"

type User struct {
	ID         int64  `db:"id"`
	TelegramID int64  `db:"telegram_id"`
	Username   string `db:"username"`
}

type Habit struct {
	ID         int64     `db:"id"`
	TelegramID int64     `db:"telegram_id"`
	Text       string    `db:"text"`
	Series     string    `db:"series"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}
