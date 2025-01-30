package main

import (
	"PB/internal/bot"
	"PB/internal/clients"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	telegramBot, err := bot.NewBot(os.Getenv("BOT_TOKEN"), os.Getenv("API_BASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	userClient := clients.NewUserClient(apiBaseURL)

	// Создаём обработчик бота
	botHandler := handlers.NewBotHandler(telegramBot, userClient)

	// Настраиваем получение обновлений (исправленный код!)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := telegramBot.GetUpdates(u)
	if err != nil {
		log.Fatal("❌ Ошибка получения обновлений:", err)
	}

	log.Println("🚀 Бот успешно запущен и ждёт команды!")

	// Обрабатываем входящие сообщения
	for _, update := range updates {
		if update.Message != nil {
			switch update.Message.Text {
			case "/start":
				botHandler.HandleStart(&update)
			}
		}
	}
}
