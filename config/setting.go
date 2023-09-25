package config

import (
	"fmt"

	"github.com/caarlos0/env/v9"
	"github.com/joho/godotenv"
)

type EnvConfig struct {
	Port         int    `env:"PORT,required" envDefault:"8080"`
	Environment  string `env:"ENVIRONMENT,required" envDefault:"development"`
	Version      string `env:"VERSION" envDefault:"0.0.1"`
	CookieSecret string `env:"COOKIE_SECRET,required"`
	DSN          string `env:"DATABASE_URL,required"`
}

func LoadEnvConfig() (EnvConfig, error) {
	err := godotenv.Load()
	if err != nil {
		return EnvConfig{}, fmt.Errorf("failed to load env: %w", err)
	}
	cfg := EnvConfig{}
	if err := env.Parse(&cfg); err != nil {
		return EnvConfig{}, fmt.Errorf("failed to parse env: %w", err)
	}
	return cfg, nil
}

func (conf EnvConfig) IsProd() bool {
	return conf.Environment == "production"
}

func (conf EnvConfig) IsDev() bool {
	return conf.Environment == "development"
}
