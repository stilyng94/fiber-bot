package config

import (
	"fmt"
	"strings"

	"github.com/caarlos0/env/v9"
	"github.com/joho/godotenv"
)

type DatabaseEngineType string

const (
	sqliteEngine   DatabaseEngineType = "sqlite"
	postgresEngine DatabaseEngineType = "postgres"
)

type EnvConfig struct {
	Port                   int    `env:"PORT,required" envDefault:"8080"`
	Environment            string `env:"ENVIRONMENT,required" envDefault:"development"`
	Version                string `env:"VERSION" envDefault:"0.0.1"`
	CookieSecret           string `env:"COOKIE_SECRET,required"`
	DSN                    string `env:"DATABASE_URL,required" envDefault:"file::memory:?cache=shared&_pragma=foreign_keys(1)"`
	TelegramBotToken       string `env:"TELEGRAM_BOT_TOKEN,required"`
	AppUrl                 string `env:"APP_URL,required"`
	TelegramStripeToken    string `env:"TELEGRAM_STRIPE_TOKEN,required"`
	Domain                 string `env:"DOMAIN,required"`
	AllowedTelegramAdmins  []string
	UnparsedTelegramAdmins string             `env:"ALLOWED_TELEGRAM_ADMINS,required"`
	CloudinaryUrl          string             `env:"CLOUDINARY_URL,required"`
	DatabaseEngine         DatabaseEngineType `env:"DATABASE_ENGINE,required" envDefault:"sqlite"`
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
	cfg.AllowedTelegramAdmins = strings.Split(cfg.UnparsedTelegramAdmins, ",")
	return cfg, nil
}

func (conf EnvConfig) IsProd() bool {
	return conf.Environment == "production"
}

func (conf EnvConfig) IsDev() bool {
	return conf.Environment == "development"
}
