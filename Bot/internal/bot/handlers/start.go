package handlers

import (
	"PB/constants"
	"PB/internal/clients"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleMessage(
	message *tgbotapi.Message,
	apiUserClient *clients.UserClient,
	bot *tgbotapi.BotAPI,
) {
	if message.Command() == "start" {
		handleStart(message, apiUserClient, bot)
	}
}

func handleStart(
	message *tgbotapi.Message,
	apiUserClient *clients.UserClient,
	bot *tgbotapi.BotAPI,
) {
	telegramID := message.From.ID
	username := message.From.UserName

	user, err := apiUserClient.GetUser(telegramID)
	fmt.Print(telegramID)
	fmt.Print(user)
	if err != nil {
		newUser := &clients.User{
			TelegramID: telegramID,
			Username:   username,
			RoleID:     1,
		}

		err := apiUserClient.CreateUser(newUser)
		if err != nil {
			bot.Send(tgbotapi.NewMessage(message.Chat.ID, constants.ErrRegistrationMsg))
			return
		}

		bot.Send(tgbotapi.NewMessage(message.Chat.ID, constants.WelcomeMsg))
		return
	}

	var keyboard tgbotapi.ReplyKeyboardMarkup

	switch user.RoleID {
	case 1:
		keyboard = tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("📊 Статистика"),
				tgbotapi.NewKeyboardButton("⚙ Управление пользователями"),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("🔧 Настройки"),
			),
		)
		msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("Добро пожаловать, %s!\nВаша роль: %v", user.Username, user.Role))
		msg.ReplyMarkup = keyboard
	case 2:
		bot.Send(tgbotapi.NewMessage(message.Chat.ID, "2"))
	}
}
