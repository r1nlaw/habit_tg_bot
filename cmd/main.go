package main

import (
	"context"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"

	"habit_bot/internal/config"
	"habit_bot/internal/handler"
	"habit_bot/internal/repository"
	"habit_bot/internal/service"
)

func main() {
	_ = godotenv.Load("./.env")
	ctx := context.Background()

	cfg := config.Load()

	db, err := repository.NewPostgresDB(cfg.DatabaseConfig)
	if err != nil {
		log.Fatal("failed to connect to DB: ", err)
	}

	bot, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		log.Fatal(err)
	}
	bot.Debug = true

	userRepo := repository.NewUserPostgres(ctx, db)
	habitRepo := repository.NewHabitPostgres(ctx, db)

	userService := service.NewUserService(ctx, userRepo)
	habitService := service.NewHabitService(habitRepo, ctx)

	tgHandler := handler.NewTelegramHandler(bot, userService, habitService)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	tgHandler.HandleUpdates(updates)
}
