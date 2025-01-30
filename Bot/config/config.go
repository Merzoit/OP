package config

import (
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	TelegramToken string
	APIBaseURL    string
}

func LoadConfig(path string) (*Config, error) {
	err := godotenv.Load(path)
	if err != nil {
		return nil, err
	}

	return &Config{
		TelegramToken: os.Getenv("TELEGRAM_TOKEN"),
		APIBaseURL:    os.Getenv("API_BASE_URL"),
	}, nil
}
