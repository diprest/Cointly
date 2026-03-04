package config

import (
	"encoding/json"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Addr           string
	DBUrl          string
	RedisAddr      string
	RedisPassword  string
	UpdateInterval time.Duration
	KafkaBrokers   []string
	KafkaTopic     string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		slog.Warn("No .env file found, relying on system environment variables")
	}

	dsn := "host=" + getEnv("DB_HOST", "localhost") +
		" user=" + getEnv("DB_USER", "postgres") +
		" password=" + getEnv("DB_PASSWORD", "postgres") +
		" dbname=" + getEnv("DB_NAME", "market_db") +
		" port=" + getEnv("DB_PORT", "5432") +
		" sslmode=disable"

	intervalMs, _ := strconv.Atoi(getEnv("UPDATE_INTERVAL_MS", "500"))

	if intervalMs <= 0 {
		intervalMs = 500
	}

	return &Config{
		Addr:          ":" + getEnv("PORT", "8001"),
		DBUrl:         dsn,
		RedisAddr:     getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),

		UpdateInterval: time.Duration(intervalMs) * time.Millisecond,

		KafkaBrokers: []string{getEnv("KAFKA_BROKERS", "localhost:9092")},
		KafkaTopic:   getEnv("KAFKA_TOPIC", "market_prices"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

type SymbolConfigItem struct {
	Symbol string `json:"symbol"`
	Name   string `json:"name"`
}

type SymbolsConfig struct {
	Symbols []SymbolConfigItem `json:"symbols"`
}

func LoadSymbols(path string) ([]SymbolConfigItem, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cfg SymbolsConfig
	if err := json.NewDecoder(file).Decode(&cfg); err != nil {
		return nil, err
	}
	return cfg.Symbols, nil
}
