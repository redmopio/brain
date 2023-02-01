package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

type Config struct {
	DatabaseURL string
	OpenAIKey   string
}

func NewLoadedConfig() (*Config, error) {
	var c Config
	err := envconfig.Process("brain", &c)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &c, nil
}
