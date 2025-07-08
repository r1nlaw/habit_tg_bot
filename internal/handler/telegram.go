package handler

import (
	"context"
	"fmt"
	"habit_bot/internal/models"
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"habit_bot/internal/service"
)

type TelegramHandler struct {
	UserService  *service.UserService
	HabitService *service.HabitService
	Bot          *tgbotapi.BotAPI
}

func NewTelegramHandler(bot *tgbotapi.BotAPI, userService *service.UserService, habitService *service.HabitService) *TelegramHandler {
	return &TelegramHandler{
		UserService:  userService,
		HabitService: habitService,
		Bot:          bot,
	}
}

func (h *TelegramHandler) HandleUpdates(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.Message == nil {
			continue
		}

		user, err := h.UserService.RegisterUser(update.Message.From.ID, update.Message.From.UserName)
		if err != nil {
			log.Println("Ошибка регистрации пользователя:", err)
			continue
		}

		text := update.Message.Text
		chatID := update.Message.Chat.ID
		telegramID := update.Message.From.ID

		if strings.HasPrefix(text, "/add ") {
			h.handleAdd(update.Message, telegramID)
		} else if text == "/list" {
			h.handleList(update.Message, telegramID)
		} else if strings.HasPrefix(text, "/delete ") {
			h.handleDelete(update.Message, telegramID)
		} else if strings.HasPrefix(text, "/update ") {
			h.handleUpdate(update.Message, telegramID)
		} else {
			msg := tgbotapi.NewMessage(chatID, "Привет, "+user.Username+"\n\n"+
				"Используй команды:\n"+
				"/add <текст> - добавить привычку\n"+
				"/list - показать привычки\n"+
				"/delete <id> - удалить привычку\n"+
				"/update <id> - обновить привычку")
			h.Bot.Send(msg)
		}
	}
}

func (h *TelegramHandler) handleAdd(msg *tgbotapi.Message, telegramID int64) {
	habitText := strings.TrimSpace(strings.TrimPrefix(msg.Text, "/add"))
	if habitText == "" {
		h.reply(msg.Chat.ID, "Пожалуйста, укажите текст привычки после команды /add")
		return
	}

	habit := models.Habit{
		TelegramID: telegramID,
		Text:       habitText,
	}

	err := h.HabitService.AddHabits(context.Background(), habit)
	if err != nil {
		h.reply(msg.Chat.ID, "Ошибка при добавлении привычки: "+err.Error())
		return
	}

	h.reply(msg.Chat.ID, "Привычка успешно добавлена!")
}

func (h *TelegramHandler) handleList(msg *tgbotapi.Message, telegramID int64) {
	habits, err := h.HabitService.GetHabits(context.Background(), telegramID)
	if err != nil {
		h.reply(msg.Chat.ID, "Ошибка при получении привычек: "+err.Error())
		return
	}
	if len(habits) == 0 {
		h.reply(msg.Chat.ID, "У вас нет привычек.")
		return
	}

	response := "Ваши привычки:\n"
	for _, habit := range habits {
		response += fmt.Sprintf("ID: %d | %s | Серия: %s\n", habit.ID, habit.Text, habit.Series)
	}

	h.reply(msg.Chat.ID, response)
}

func (h *TelegramHandler) handleDelete(msg *tgbotapi.Message, telegramID int64) {
	idStr := strings.TrimSpace(strings.TrimPrefix(msg.Text, "/delete"))
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.reply(msg.Chat.ID, "Некорректный ID привычки.")
		return
	}

	err = h.HabitService.DeleteHabits(context.Background(), id)
	if err != nil {
		h.reply(msg.Chat.ID, "Ошибка при удалении привычки: "+err.Error())
		return
	}

	h.reply(msg.Chat.ID, "Привычка успешно удалена.")
}

func (h *TelegramHandler) handleUpdate(msg *tgbotapi.Message, telegramID int64) {
	idStr := strings.TrimSpace(strings.TrimPrefix(msg.Text, "/update"))
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.reply(msg.Chat.ID, "Некорректный ID привычки.")
		return
	}

	err = h.HabitService.UpdateHabit(context.Background(), id)
	if err != nil {
		h.reply(msg.Chat.ID, "Ошибка при обновлении привычки: "+err.Error())
		return
	}

	h.reply(msg.Chat.ID, "Привычка успешно обновлена.")
}

func (h *TelegramHandler) reply(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	if _, err := h.Bot.Send(msg); err != nil {
		log.Println("Ошибка отправки сообщения:", err)
	}
}
