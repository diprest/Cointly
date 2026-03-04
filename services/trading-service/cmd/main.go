package main

import (
	"context"
	"log"
	"strings"
	"time"

	"trading-service/internal/broker"
	"trading-service/internal/clients"
	"trading-service/internal/config"
	"trading-service/internal/models"
	"trading-service/internal/service"
	"trading-service/internal/storage"
	"trading-service/internal/transport/http"
	"trading-service/internal/worker"
)

func main() {
	cfg := config.Load()

	store, err := storage.NewStorage(cfg.DBUrl)
	if err != nil {
		log.Fatalf("DB Error: %v", err)
	}

	portfolioURL := cfg.PortfolioURL
	if portfolioURL == "" {
		portfolioURL = "http://localhost:8083"
	}
	pfClient := clients.NewPortfolioHTTP(portfolioURL)
	log.Printf("🔌 Connected to Portfolio Service at %s", portfolioURL)

	marketDataURL := cfg.MarketDataURL
	if marketDataURL == "" {
		marketDataURL = "http://localhost:8001"
	}
	mdClient := clients.NewMarketDataHTTP(marketDataURL)
	log.Printf("🔌 Connected to Market Data Service at %s", marketDataURL)

	priceChannel := make(chan models.PriceUpdate, 100)

	brokers := strings.Split(cfg.KafkaBroker, ",")
	consumer, err := broker.NewKafkaConsumer(brokers, "market_prices")
	if err != nil {
		log.Printf("Kafka init failed: %v. Running without live updates.", err)
	} else {
		consumer.Subscribe(priceChannel)
	}

	matcher := worker.NewMatcher(store, pfClient, mdClient, priceChannel)
	go matcher.Start()
	go matcher.StartPolling(context.Background(), 5*time.Second)
	svc := service.NewTradingService(store, pfClient, mdClient)

	handler := http.NewHandler(svc)

	r := handler.InitRoutes()

	log.Printf("Starting server on %s", cfg.Port)
	if err := r.Run(cfg.Port); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
