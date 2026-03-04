package service

import (
	"context"
	"market-data-service/internal/models"
	"time"
)

type MarketService struct {
	symbolRepo SymbolRepository
	priceRepo  PriceRepository
}

func NewMarketService(s SymbolRepository, p PriceRepository) *MarketService {
	return &MarketService{
		symbolRepo: s,
		priceRepo:  p,
	}
}

func (s *MarketService) GetAllSymbols(ctx context.Context) ([]models.CoinInfo, error) {
	symbols, err := s.symbolRepo.GetActiveSymbols()
	if err != nil {
		return nil, err
	}

	var coinInfos []models.CoinInfo
	for _, sym := range symbols {
		price, _ := s.priceRepo.GetPrice(ctx, sym.Symbol)

		pnl := 0.0
		if price > 0 {
			if cachedPnL, err := s.priceRepo.GetCachedPnL(ctx, sym.Symbol); err == nil {
				pnl = cachedPnL
			} else {
				if oldPrice, err := s.priceRepo.GetOldestPrice(ctx, sym.Symbol); err == nil && oldPrice > 0 {
					pnl = ((price - oldPrice) / oldPrice) * 100
					_ = s.priceRepo.CachePnL(ctx, sym.Symbol, pnl, 5*time.Minute)
				}
			}
		}

		coinInfos = append(coinInfos, models.CoinInfo{
			Symbol: sym.Symbol,
			Name:   sym.Name,
			Price:  price,
			PnL:    pnl,
		})
	}
	return coinInfos, nil
}

func (s *MarketService) GetTicker(ctx context.Context, symbol string) (models.TickerPrice, error) {
	price, err := s.priceRepo.GetPrice(ctx, symbol)
	if err != nil {
		return models.TickerPrice{}, err
	}
	return models.TickerPrice{
		Symbol: symbol,
		Price:  price,
	}, nil
}
