package repository

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"habit_bot/internal/models"
)

type habitDB struct {
	ctx      context.Context
	postgres *sqlx.DB
}

func NewHabitPostgres(ctx context.Context, db *sqlx.DB) *habitDB {
	return &habitDB{ctx: ctx, postgres: db}
}

func (h *habitDB) AddHabits(ctx context.Context, habit models.Habit) error {
	query := `INSERT INTO habits (telegram_id, text) VALUES ($1, $2)`
	_, err := h.postgres.ExecContext(ctx, query, habit.TelegramID, habit.Text)
	return err
}

func (h *habitDB) DeleteHabits(ctx context.Context, id int64) error {
	query := `DELETE FROM habits WHERE id = $1`
	result, err := h.postgres.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	count, _ := result.RowsAffected()
	if count == 0 {
		return fmt.Errorf("no habit with id %d", id)
	}
	return nil
}

func (h *habitDB) GetHabits(ctx context.Context, telegramID int64) ([]models.Habit, error) {
	var habits []models.Habit
	query := `SELECT * FROM habits WHERE telegram_id = $1 ORDER BY created_at DESC`
	err := h.postgres.SelectContext(ctx, &habits, query, telegramID)
	return habits, err
}

func (h *habitDB) UpdateHabit(ctx context.Context, id int64) error {
	query := `UPDATE habits SET series = series + 1, updated_at = NOW() WHERE id = $1`
	result, err := h.postgres.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	count, _ := result.RowsAffected()
	if count == 0 {
		return fmt.Errorf("no habit with id %d", id)
	}
	return nil
}
