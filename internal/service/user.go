package service

import (
	"context"
	"habit_bot/internal/models"
	"habit_bot/internal/repository"
)

type UserService struct {
	repo repository.UserRepository
	ctx  context.Context
}

func NewUserService(ctx context.Context, repo repository.UserRepository) *UserService {
	return &UserService{repo: repo, ctx: ctx}
}

func (u *UserService) RegisterUser(telegramID int64, username string) (models.User, error) {
	return u.repo.GetOrCreate(u.ctx, telegramID, username)
}
