package router

import (
	"net/http"

	"api-gateway/internal/config"
	"api-gateway/internal/middleware"
	"api-gateway/internal/proxy"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func NewRouter(cfg *config.Config) http.Handler {
	r := chi.NewRouter()

	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	authMiddleware := middleware.NewAuthMiddleware(cfg.JWTSecret)

	authProxy := proxy.NewReverseProxy(cfg.AuthServiceURL)
	marketProxy := proxy.NewReverseProxy(cfg.MarketDataURL)
	tradingProxy := proxy.NewReverseProxy(cfg.TradingURL)
	portfolioProxy := proxy.NewReverseProxy(cfg.PortfolioURL)
	betsProxy := proxy.NewReverseProxy(cfg.BetsURL)

	r.Mount("/api/v1/auth", authProxy)

	r.Group(func(r chi.Router) {
		r.Use(authMiddleware.VerifyToken)

		r.Mount("/api/v1/market", marketProxy)

		r.Mount("/api/v1/trading", tradingProxy)
		r.Mount("/api/v1/portfolio", portfolioProxy)
		r.Mount("/api/v1/bets", betsProxy)
	})

	return r
}
