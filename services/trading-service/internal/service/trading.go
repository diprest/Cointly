package service

import (
	"context"
	"errors"
	"fmt"
	"trading-service/internal/clients"
	"trading-service/internal/models"
	"trading-service/internal/storage"

	"github.com/shopspring/decimal"
)

type TradingService struct {
	repo       storage.OrderRepository
	portfolio  clients.PortfolioClient
	marketData clients.MarketDataClient
}

func NewTradingService(repo storage.OrderRepository, pf clients.PortfolioClient, md clients.MarketDataClient) *TradingService {
	return &TradingService{
		repo:       repo,
		portfolio:  pf,
		marketData: md,
	}
}

func (s *TradingService) CreateOrder(o *models.Order) error {
	o.QuoteAmount = o.QuoteAmount.Abs()

	if o.Side == "BUY" && o.Type == "MARKET" {
		if o.QuoteAmount.GreaterThan(decimal.Zero) {
			o.Amount = decimal.Zero
		} else if o.Amount.GreaterThan(decimal.Zero) {
			if o.Price.LessThanOrEqual(decimal.Zero) {
				return errors.New("for MARKET BUY by quantity, provide estimated price")
			}
		} else {
			return errors.New("provide either amount or quote_amount")
		}
	} else {
		if o.Amount.LessThanOrEqual(decimal.Zero) {
			return errors.New("amount must be positive")
		}
		if o.Type == "LIMIT" && o.Price.LessThanOrEqual(decimal.Zero) {
			return errors.New("price must be positive")
		}
	}

	o.Status = "NEW"

	shouldExecuteImmediately := false
	if o.Type == "MARKET" {
		shouldExecuteImmediately = true
	} else if o.Type == "LIMIT" {
		price, err := s.marketData.GetPrice(o.Symbol)
		if err == nil && price.GreaterThan(decimal.Zero) {
			if o.Side == "BUY" && o.Price.GreaterThanOrEqual(price) {
				shouldExecuteImmediately = true
			} else if o.Side == "SELL" && o.Price.LessThanOrEqual(price) {
				shouldExecuteImmediately = true
			}
		}
	}

	if shouldExecuteImmediately {
		price, err := s.marketData.GetPrice(o.Symbol)
		if err == nil && price.GreaterThan(decimal.Zero) {
			var cost, amount decimal.Decimal

			execPrice := price

			if o.Side == "BUY" {
				if o.QuoteAmount.GreaterThan(decimal.Zero) {
					cost = o.QuoteAmount
					amount = o.QuoteAmount.Div(execPrice)
					o.Amount = amount
				} else {
					amount = o.Amount
					cost = o.Amount.Mul(execPrice)
				}
			} else {
				amount = o.Amount
				cost = o.Amount.Mul(execPrice)
			}

			asset := o.Symbol
			if len(o.Symbol) > 4 && o.Symbol[len(o.Symbol)-4:] == "USDT" {
				asset = o.Symbol[:len(o.Symbol)-4]
			}

			err = s.portfolio.TransferFunds(context.Background(), o.UserID, asset, amount, cost, o.Side)
			if err == nil {
				o.Status = "FILLED"
				o.Price = execPrice
			} else {
				return fmt.Errorf("failed to execute order: %w", err)
			}
		}
	}

	return s.repo.CreateOrder(o)
}

func (s *TradingService) CancelOrder(ctx context.Context, orderID uint, userID int64) error {
	order, err := s.repo.GetOrder(orderID)
	if err != nil {
		return err
	}

	if order.UserID != userID {
		return errors.New("order belongs to another user")
	}

	if order.Status != "NEW" {
		return errors.New("can only cancel NEW orders")
	}

	return s.repo.UpdateOrderStatus(orderID, "CANCELLED")
}

func (s *TradingService) GetUserOrders(userID int64) ([]models.Order, error) {
	return s.repo.GetUserOrders(userID)
}

func (s *TradingService) GetUserOrdersWithCtx(ctx context.Context, userID int64) ([]models.Order, error) {
	return s.GetUserOrders(userID)
}

func (s *TradingService) ResetOrders(ctx context.Context, userID int64) error {
	return s.repo.ResetOrders(userID)
}
