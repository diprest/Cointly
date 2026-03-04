package main

import (
	"log"
	"net/http"
	"portfolio-service/internal/config"
	"portfolio-service/internal/models"
	"portfolio-service/internal/service"
	"portfolio-service/internal/storage"
	h "portfolio-service/internal/transport/http"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg := config.LoadConfig()

	dbConn, err := gorm.Open(postgres.Open(cfg.DBUrl), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	if err := dbConn.AutoMigrate(&models.Balance{}); err != nil {
		log.Fatalf("Failed to migrate DB: %v", err)
	}
	repo := storage.NewPostgresDB(dbConn)
	svc := service.NewPortfolioService(repo)
	handler := h.NewHandler(svc)
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)
	log.Printf("Portfolio Service running on %s", cfg.Addr)
	if err := http.ListenAndServe(cfg.Addr, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
