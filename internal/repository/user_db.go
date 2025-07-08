package repository

import (
	"context"
	"github.com/jmoiron/sqlx"
	"habit_bot/internal/models"
)

type userDB struct {
	ctx      context.Context
	postgres *sqlx.DB
}

func NewUserPostgres(ctx context.Context, db *sqlx.DB) *userDB {
	return &userDB{ctx: ctx, postgres: db}
}

func (u *userDB) GetOrCreate(ctx context.Context, telegramID int64, username string) (models.User, error) {
	var user models.User
	err := u.postgres.GetContext(ctx, &user, "SELECT * FROM users WHERE telegram_id=$1", telegramID)
	if err == nil {
		return user, nil
	}
	query := `INSERT INTO users (telegram_id, username) VALUES ($1, $2) RETURNING *`
	err = u.postgres.GetContext(ctx, &user, query, telegramID, username)
	return user, err
}
