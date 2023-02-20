package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

type Config struct {
	OpenAIKey            string `envconfig:"openai_key" required:"true"`
	Name                 string `envconfig:"name" default:"AI"`
	DatabaseURL          string `envconfig:"database_url" default:"sqlite3://brain.db"`
	WhatsAppDisabled     bool   `envconfig:"whatsapp_disabled" default:"false"`
	WhatsAppDatabaseName string `envconfig:"whatsapp_database_name" default:"whatsapp-store.db"`
	TelegramAPIKey       string `envconfig:"telegram_api_key" default:""`
}

func NewLoadedConfig() (*Config, error) {
	godotenv.Load()

	var c Config
	err := envconfig.Process("brain", &c)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &c, nil
}
