package service

import (
	"context"
	"fmt"
	"time"
	"trading-service/internal/models"

	"github.com/shopspring/decimal"
)

func (s *TradingService) StartMatcher(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	fmt.Println("Matcher started with interval:", interval)

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Matcher stopped")
			return
		case <-ticker.C:
			s.matchOrders()
		}
	}
}

func (s *TradingService) matchOrders() {
	orders, err := s.repo.GetActiveOrders()
	if err != nil {
		fmt.Println("Matcher: failed to get active orders:", err)
		return
	}

	if len(orders) == 0 {
		return
	}

	ordersBySymbol := make(map[string][]*models.Order)
	for i := range orders {
		ordersBySymbol[orders[i].Symbol] = append(ordersBySymbol[orders[i].Symbol], &orders[i])
	}

	for symbol, symbolOrders := range ordersBySymbol {
		price, err := s.marketData.GetPrice(symbol)
		if err != nil {
			fmt.Printf("Matcher: failed to get price for %s: %v\n", symbol, err)
			continue
		}

		if price.LessThanOrEqual(decimal.Zero) {
			continue
		}

		for _, o := range symbolOrders {
			shouldExecute := false
			if o.Type == "LIMIT" {
				if o.Side == "BUY" && o.Price.GreaterThanOrEqual(price) {
					shouldExecute = true
				} else if o.Side == "SELL" && o.Price.LessThanOrEqual(price) {
					shouldExecute = true
				}
			}

			if shouldExecute {
				var cost, amount decimal.Decimal
				execPrice := price

				if o.Side == "BUY" {
					amount = o.Amount
					cost = o.Amount.Mul(execPrice)
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

					if updateErr := s.repo.UpdateOrderStatus(o.ID, "FILLED"); updateErr != nil {
						fmt.Printf("Matcher: failed to update order status %d: %v\n", o.ID, updateErr)
					} else {
						fmt.Printf("Matcher: executed order %d (%s %s) at %s\n", o.ID, o.Side, o.Symbol, execPrice)
					}
				} else {
					fmt.Printf("Matcher: failed to execute order %d: %v\n", o.ID, err)
				}
			}
		}
	}
}
