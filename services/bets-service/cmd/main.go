package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"bets-service/internal/clients"
	"bets-service/internal/config"
	"bets-service/internal/service"
	"bets-service/internal/storage"
	h "bets-service/internal/transport/http"
	"bets-service/internal/worker"
)

func main() {
	cfg := config.Load()
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	slog.Info("Starting Bets Service", "port", cfg.Port)

	store, err := storage.NewStorage(cfg.DBUrl)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer store.DB.Close()

	if err := store.RunMigrations(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	pfClient := clients.NewPortfolioClient(cfg.PortfolioURL)
	mdClient := clients.NewMarketDataClient(cfg.MarketDataURL)

	svc := service.NewBetsService(store, pfClient, mdClient)

	resolver := worker.NewResolver(svc)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go resolver.Start(ctx)

	handler := h.NewHandler(svc)
	router := handler.InitRoutes()

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("Shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	slog.Info("Server exiting")
}
