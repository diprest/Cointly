package main

import (
	"auth-service/internal/config"
	"auth-service/internal/db"
	"auth-service/internal/repository"
	"auth-service/internal/service"
	transport "auth-service/internal/transport/http"
	"log"
	"net/http"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	log.Printf("Config loaded: Port=%s, DB_URL=%s", cfg.Port, cfg.DatabaseURL)

	dbPool, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer dbPool.Close()

	if err := db.RunMigrations(dbPool); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	userRepo := repository.NewUserRepository(dbPool)
	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	handler := transport.NewHandler(authService)

	r := handler.InitRoutes()

	log.Printf("Gateway running on port %s", cfg.Port)
	if err := http.ListenAndServe("0.0.0.0:"+cfg.Port, r); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
