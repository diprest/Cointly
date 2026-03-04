package config

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUrl         string
	KafkaBroker   string
	Port          string
	PortfolioURL  string
	MarketDataURL string
}

func Load() *Config {
	_ = godotenv.Load()

	kafka := os.Getenv("KAFKA_BROKERS")
	if kafka == "" {
		kafka = os.Getenv("KAFKA_BROKER")
		if kafka == "" {
			kafka = "kafka:29092"
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = ":8082"
	}
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}

	cfg := &Config{
		DBUrl:         os.Getenv("DATABASE_URL"),
		KafkaBroker:   kafka,
		Port:          port,
		PortfolioURL:  os.Getenv("PORTFOLIO_SERVICE_URL"),
		MarketDataURL: os.Getenv("MARKET_DATA_URL"),
	}

	return cfg
}
