package service

import (
	"context"
	"market-data-service/internal/models"
	"time"
)

type SymbolRepository interface {
	GetActiveSymbols() ([]models.Symbol, error)
}

type PriceRepository interface {
	GetPrice(ctx context.Context, symbol string) (float64, error)
	SetPrice(ctx context.Context, symbol string, price float64) error
	GetOldestPrice(ctx context.Context, symbol string) (float64, error)
	CachePnL(ctx context.Context, symbol string, pnl float64, ttl time.Duration) error
	GetCachedPnL(ctx context.Context, symbol string) (float64, error)
}
