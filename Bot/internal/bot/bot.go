package bot

import (
	"PB/internal/bot/handlers"
	"PB/internal/clients"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	botAPI        *tgbotapi.BotAPI
	apiUserClient *clients.UserClient
}

func NewBot(token, apiBaseURL string) (*Bot, error) {
	botAPI, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	apiClient := clients.NewUserClient(apiBaseURL)

	return &Bot{
		botAPI:        botAPI,
		apiUserClient: apiClient,
	}, nil
}

func (b *Bot) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.botAPI.GetUpdatesChan(u)
	for update := range updates {
		if update.Message != nil {
			handlers.HandleMessage(update.Message, b.apiUserClient, b.botAPI)
		}
	}
}
