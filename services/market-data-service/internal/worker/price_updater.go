package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"market-data-service/internal/models"
	"net/http"
	"time"
)

type SymbolProvider interface {
	GetActiveSymbols() ([]models.Symbol, error)
}

type PriceSaver interface {
	SetPrice(ctx context.Context, symbol string, price float64) error
}

type HistorySaver interface {
	AddPriceHistory(ctx context.Context, symbol string, price float64, timestamp int64) error
	TrimHistory(ctx context.Context, symbol string, retention int64) error
}

type EventPublisher interface {
	Publish(symbol string, price float64) error
}

type BinanceTicker struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

type PriceUpdater struct {
	pg          SymbolProvider
	redis       PriceSaver
	broker      EventPublisher
	interval    time.Duration
	client      *http.Client
	apiUrl      string
	lastHistory map[string]time.Time
}

func NewPriceUpdater(pg SymbolProvider, redis PriceSaver, broker EventPublisher, interval time.Duration) *PriceUpdater {
	return &PriceUpdater{
		pg:          pg,
		redis:       redis,
		broker:      broker,
		interval:    interval,
		client:      &http.Client{Timeout: 5 * time.Second},
		apiUrl:      "https://api.binance.com/api/v3/ticker/price",
		lastHistory: make(map[string]time.Time),
	}
}

func (w *PriceUpdater) SetAPIURL(url string) {
	w.apiUrl = url
}

func (w *PriceUpdater) Start(ctx context.Context) {
	slog.Info("Price Updater Worker started", "interval", w.interval)
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	w.updatePrices(ctx)

	for {
		select {
		case <-ctx.Done():
			slog.Info("Price Updater Worker stopping...")
			return
		case <-ticker.C:
			w.updatePrices(ctx)
		}
	}
}

func (w *PriceUpdater) updatePrices(ctx context.Context) {
	symbols, err := w.pg.GetActiveSymbols()
	if err != nil {
		slog.Error("Failed to fetch symbols", "error", err)
		return
	}

	if len(symbols) == 0 {
		return
	}

	resp, err := w.client.Get(w.apiUrl)
	if err != nil {
		slog.Error("Failed to fetch prices from External API", "error", err)
		return
	}
	defer resp.Body.Close()

	var binanceTickers []BinanceTicker
	if err := json.NewDecoder(resp.Body).Decode(&binanceTickers); err != nil {
		slog.Error("Failed to decode response", "error", err)
		return
	}

	priceMap := make(map[string]float64, len(binanceTickers))
	for _, t := range binanceTickers {
		var p float64
		if _, err := fmt.Sscanf(t.Price, "%f", &p); err == nil {
			priceMap[t.Symbol] = p
		}
	}

	updatedCount := 0
	now := time.Now()
	historyInterval := 5 * time.Minute
	retention := now.Add(-24*time.Hour - 10*time.Minute).Unix()

	for _, s := range symbols {
		if price, ok := priceMap[s.Symbol]; ok {
			if err := w.redis.SetPrice(ctx, s.Symbol, price); err != nil {
				slog.Error("Redis save error", "symbol", s.Symbol, "error", err)
			} else {
				updatedCount++
			}

			if w.broker != nil {
				if err := w.broker.Publish(s.Symbol, price); err != nil {
					slog.Error("Kafka publish error", "symbol", s.Symbol, "error", err)
				}
			}

			lastUpdate, exists := w.lastHistory[s.Symbol]
			if !exists || now.Sub(lastUpdate) >= historyInterval {
				if saver, ok := w.redis.(HistorySaver); ok {
					if err := saver.AddPriceHistory(ctx, s.Symbol, price, now.Unix()); err != nil {
						slog.Error("Failed to add price history", "symbol", s.Symbol, "error", err)
					}
					if err := saver.TrimHistory(ctx, s.Symbol, retention); err != nil {
						slog.Error("Failed to trim history", "symbol", s.Symbol, "error", err)
					}
					w.lastHistory[s.Symbol] = now
				}
			}
		}
	}

	slog.Debug("Market data updated", "count", updatedCount)
}
