package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"api-gateway/internal/config"
	"api-gateway/internal/router"
)

func main() {
	cfg := config.Load()
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	slog.Info("Starting API Gateway", "port", cfg.Port)

	r := router.NewRouter(cfg)

	log.Printf("Gateway running on port %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
