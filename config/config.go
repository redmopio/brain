package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

type Config struct {
	DatabaseURL string `envconfig:"database_url" required:"true"`
	OpenAIKey   string `envconfig:"openai_key" required:"true"`
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
