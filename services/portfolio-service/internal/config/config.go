package config

import (
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
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		slog.Warn("No .env file found, relying on system environment variables")
	}

	dsn := "host=" + getEnv("DB_HOST", "localhost") +
		" user=" + getEnv("DB_USER", "postgres") +
		" password=" + getEnv("DB_PASSWORD", "postgres") +
		" dbname=" + getEnv("DB_NAME", "portfolio_db") +
		" port=" + getEnv("DB_PORT", "5432") +
		" sslmode=disable"

	intervalSec, _ := strconv.Atoi(getEnv("UPDATE_INTERVAL_SEC", "2"))

	return &Config{
		Addr:           ":" + getEnv("PORT", "8001"),
		DBUrl:          dsn,
		RedisAddr:      getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword:  getEnv("REDIS_PASSWORD", ""),
		UpdateInterval: time.Duration(intervalSec) * time.Second,
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
