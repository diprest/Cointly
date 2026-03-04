package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port          string
	DBUrl         string
	RedisAddr     string
	PortfolioURL  string
	MarketDataURL string
	JWTSecret     string
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		Port:          getEnv("PORT", "8084"),
		DBUrl:         getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/bets_db?sslmode=disable"),
		RedisAddr:     getEnv("REDIS_ADDR", "localhost:6379"),
		PortfolioURL:  getEnv("PORTFOLIO_URL", "http://localhost:8083"),
		MarketDataURL: getEnv("MARKET_DATA_URL", "http://localhost:8001"),
		JWTSecret:     getEnv("JWT_SECRET", "secret"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
