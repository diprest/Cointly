package main

import (
	"context"
	"log/slog"
	"market-data-service/internal/broker"
	"market-data-service/internal/config"
	"market-data-service/internal/service"
	"market-data-service/internal/storage"
	httpHandler "market-data-service/internal/transport/http"
	"market-data-service/internal/worker"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg := config.LoadConfig()

	symbols, err := config.LoadSymbols("configs/crypto_symbols.json")
	if err != nil {
		slog.Warn("Failed to load symbols from config, using defaults", "error", err)
		symbols = []config.SymbolConfigItem{
			{Symbol: "BTCUSDT", Name: "Bitcoin"},
			{Symbol: "ETHUSDT", Name: "Ethereum"},
		}
	}

	pgDB, err := storage.NewPostgresDB(cfg.DBUrl, symbols)
	if err != nil {
		slog.Error("Failed to connect to Postgres", "error", err)
		os.Exit(1)
	}
	slog.Info("Connected to Postgres")

	redisClient, err := storage.NewRedisClient(cfg.RedisAddr, cfg.RedisPassword)
	if err != nil {
		slog.Error("Failed to connect to Redis", "error", err)
		os.Exit(1)
	}
	slog.Info("Connected to Redis")

	var kafkaProducer *broker.KafkaProducer
	for i := 0; i < 10; i++ {
		kafkaProducer, err = broker.NewKafkaProducer(cfg.KafkaBrokers, cfg.KafkaTopic)
		if err == nil {
			break
		}
		slog.Warn("Failed to connect to Kafka, retrying...", "attempt", i+1, "error", err)
		time.Sleep(5 * time.Second)
	}
	if err != nil {
		slog.Error("Failed to connect to Kafka after retries", "brokers", cfg.KafkaBrokers, "error", err)
	} else {
		defer kafkaProducer.Close()
		slog.Info("Connected to Kafka", "brokers", cfg.KafkaBrokers)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	priceUpdater := worker.NewPriceUpdater(pgDB, redisClient, kafkaProducer, cfg.UpdateInterval)
	go priceUpdater.Start(ctx)

	svc := service.NewMarketService(pgDB, redisClient)
	handler := httpHandler.NewHandler(svc)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	srv := &http.Server{
		Addr:         cfg.Addr,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		slog.Info("Starting server", "addr", cfg.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server startup failed", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("Shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
	}

	slog.Info("Server exited properly")
}
