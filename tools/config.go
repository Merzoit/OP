package tools

import (
	"fmt"
	"os"

	"at/constants"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

type Config struct {
	App struct {
		Port int `yaml:"port"`
	} `yaml:"app"`
	Database struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Name     string `yaml:"dbname"`
	} `yaml:"database"`
}

func (c *Config) Validate() error {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	if c.App.Port == 0 {
		log.Error().Msg(constants.ErrConfigValidateAppPort)
		return fmt.Errorf(constants.ErrConfigValidateAppPort)
	}

	if c.Database.Host == "" || c.Database.Name == "" {
		log.Error().Msg(constants.ErrConfigValidateDbHost)
		return fmt.Errorf(constants.ErrConfigValidateDbHost)
	}

	log.Info().Msg("Configuration validation passed successfully")
	return nil
}

func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	decoder := yaml.NewDecoder(file)

	err = decoder.Decode(&config)
	return &config, err
}
