package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Port        string `env:"PORT" env-default:"8081"`
	DatabaseURL string `env:"DATABASE_URL" env-required:"true"`
	JWTSecret   string `env:"JWT_SECRET" env-required:"true"`
}

func Load() (*Config, error) {
	_ = godotenv.Load("../../.env")
	_ = godotenv.Load()

	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	return &cfg, nil
}
