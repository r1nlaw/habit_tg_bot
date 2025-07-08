package repository

import (
	"context"
	"habit_bot/internal/models"
)

type Repository struct {
	HabitRepository
	UserRepository
}
type UserRepository interface {
	GetOrCreate(ctx context.Context, telegramID int64, username string) (models.User, error)
}

type HabitRepository interface {
	AddHabits(ctx context.Context, habit models.Habit) error
	DeleteHabits(ctx context.Context, id int64) error
	GetHabits(ctx context.Context, telegramID int64) ([]models.Habit, error)
	UpdateHabit(ctx context.Context, id int64) error
}
