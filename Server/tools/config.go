package tools

import (
	"fmt"
	"os"
	"strconv"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Config struct {
	AppPort      int
	DatabaseHost string
	DatabasePort int
	DatabaseUser string
	DatabasePass string
	DatabaseName string
}

func (c *Config) Validate() error {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	if c.AppPort == 0 {
		log.Warn().Msg("App port is not set")
		return fmt.Errorf("app port is not set")
	}

	if c.DatabaseHost == "" || c.DatabaseName == "" {
		log.Warn().Msg("Database configuration is incomplete")
		return fmt.Errorf("database configuration is incomplete")
	}

	log.Info().Msg("Configuration validation passed successfully")
	return nil
}

func LoadConfig() (*Config, error) {
	appPort, err := strconv.Atoi(os.Getenv("APP_PORT"))
	if err != nil {
		log.Warn().Msg("Invalid app port")
		return nil, fmt.Errorf("invalid app port: %w", err)
	}

	dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		log.Warn().Msg("Invalid database port")
		return nil, fmt.Errorf("invalid database port: %w", err)
	}

	config := &Config{
		AppPort:      appPort,
		DatabaseHost: os.Getenv("DB_HOST"),
		DatabasePort: dbPort,
		DatabaseUser: os.Getenv("DB_USER"),
		DatabasePass: os.Getenv("DB_PASSWORD"),
		DatabaseName: os.Getenv("DB_NAME"),
	}

	return config, nil
}
