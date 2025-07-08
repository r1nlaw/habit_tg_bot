package service

import (
	"context"
	"habit_bot/internal/models"
	"habit_bot/internal/repository"
)

type HabitService struct {
	repo repository.HabitRepository
	ctx  context.Context
}

func NewHabitService(repo repository.HabitRepository, ctx context.Context) *HabitService {
	return &HabitService{repo: repo, ctx: ctx}
}

func (h *HabitService) GetHabits(ctx context.Context, telegramID int64) ([]models.Habit, error) {
	return h.repo.GetHabits(ctx, telegramID)
}
func (h *HabitService) AddHabits(ctx context.Context, habit models.Habit) error {
	return h.repo.AddHabits(ctx, habit)
}

func (h *HabitService) DeleteHabits(ctx context.Context, id int64) error {
	return h.repo.DeleteHabits(ctx, id)
}

func (h *HabitService) UpdateHabit(ctx context.Context, id int64) error {
	return h.repo.UpdateHabit(ctx, id)
}
