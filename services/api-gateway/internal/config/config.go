package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port           string
	AuthServiceURL string
	MarketDataURL  string
	TradingURL     string
	PortfolioURL   string
	BetsURL        string
	JWTSecret      string
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		Port:           getEnv("PORT", "8080"),
		AuthServiceURL: getEnv("AUTH_SERVICE_URL", "http://localhost:8081"),
		MarketDataURL:  getEnv("MARKET_DATA_URL", "http://localhost:8001"),
		TradingURL:     getEnv("TRADING_SERVICE_URL", "http://localhost:8082"),
		PortfolioURL:   getEnv("PORTFOLIO_SERVICE_URL", "http://localhost:8083"),
		BetsURL:        getEnv("BETS_SERVICE_URL", "http://localhost:8084"),
		JWTSecret:      getEnv("JWT_SECRET", "secret"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
