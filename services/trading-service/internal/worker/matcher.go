package worker

import (
	"context"
	"log"
	"time"
	"trading-service/internal/clients"
	"trading-service/internal/models"
	"trading-service/internal/storage"

	"github.com/shopspring/decimal"
)

type Matcher struct {
	store        *storage.Storage
	portfolio    clients.PortfolioClient
	marketData   clients.MarketDataClient
	priceChannel <-chan models.PriceUpdate
}

func NewMatcher(store *storage.Storage, pf clients.PortfolioClient, md clients.MarketDataClient, ch <-chan models.PriceUpdate) *Matcher {
	return &Matcher{store: store, portfolio: pf, marketData: md, priceChannel: ch}
}

func (m *Matcher) Start() {
	log.Println("⚡ Matcher Worker started (Kafka)...")
	if m.priceChannel == nil {
		log.Println("⚠️ Price channel is nil, Kafka matcher disabled")
		return
	}
	for update := range m.priceChannel {
		price, _ := decimal.NewFromString(update.Price)
		m.matchOrders(update.Symbol, price)
	}
}

func (m *Matcher) StartPolling(ctx context.Context, interval time.Duration) {
	log.Println("Matcher Polling started...")
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			m.pollAndMatch()
		}
	}
}

func (m *Matcher) pollAndMatch() {
	orders, err := m.store.GetActiveOrders()
	if err != nil {
		log.Printf("Matcher Polling DB Error: %v", err)
		return
	}
	if len(orders) == 0 {
		return
	}

	ordersBySymbol := make(map[string][]*models.Order)
	for i := range orders {
		ordersBySymbol[orders[i].Symbol] = append(ordersBySymbol[orders[i].Symbol], &orders[i])
	}

	for symbol, _ := range ordersBySymbol {
		price, err := m.marketData.GetPrice(symbol)
		if err != nil {
			log.Printf("Failed to get price for %s: %v", symbol, err)
			continue
		}
		if price.GreaterThan(decimal.Zero) {
			m.matchOrders(symbol, price)
		}
	}
}

func (m *Matcher) matchOrders(symbol string, currentPrice decimal.Decimal) {
	ctx := context.Background()

	buyQuery := `
		SELECT id, user_id, symbol, side, type, price, amount, quote_amount, status 
		FROM orders 
		WHERE symbol=$1 AND status='NEW' AND side='BUY' 
		AND ((type = 'LIMIT' AND price >= $2) OR (type = 'MARKET'))
	`
	m.processBatch(ctx, buyQuery, symbol, currentPrice)

	sellQuery := `
		SELECT id, user_id, symbol, side, type, price, amount, quote_amount, status 
		FROM orders 
		WHERE symbol=$1 AND status='NEW' AND side='SELL' 
		AND ((type = 'LIMIT' AND price <= $2) OR (type = 'MARKET'))
	`
	m.processBatch(ctx, sellQuery, symbol, currentPrice)
}

func (m *Matcher) processBatch(ctx context.Context, query string, symbol string, currentPrice decimal.Decimal) {
	rows, err := m.store.DB.Query(query, symbol, currentPrice)
	if err != nil {
		log.Printf("Matcher DB Error: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var o models.Order
		if err := rows.Scan(&o.ID, &o.UserID, &o.Symbol, &o.Side, &o.Type, &o.Price, &o.Amount, &o.QuoteAmount, &o.Status); err != nil {
			log.Printf("Scan error: %v", err)
			continue
		}
		m.executeOrder(ctx, &o, currentPrice)
	}
}

func (m *Matcher) executeOrder(ctx context.Context, o *models.Order, executionPrice decimal.Decimal) {
	log.Printf("Executing Order #%d (%s)", o.ID, o.Type)

	var cost decimal.Decimal
	var amount decimal.Decimal

	if o.QuoteAmount.GreaterThan(decimal.Zero) && o.Side == "BUY" {
		cost = o.QuoteAmount
		amount = o.QuoteAmount.Div(executionPrice)
		o.Amount = amount
	} else {
		amount = o.Amount
		cost = o.Amount.Mul(executionPrice)
	}

	asset := "BTC"
	if len(o.Symbol) > 4 && o.Symbol[len(o.Symbol)-4:] == "USDT" {
		asset = o.Symbol[:len(o.Symbol)-4]
	}

	err := m.portfolio.TransferFunds(ctx, o.UserID, asset, amount, cost, o.Side)
	if err != nil {
		log.Printf("Settle failed: %v", err)
		return
	}

	m.store.DB.Exec(`UPDATE orders SET status = 'FILLED', updated_at = $1 WHERE id = $2`, time.Now(), o.ID)
	log.Printf("Order #%d FILLED. %s %s @ %s", o.ID, amount.StringFixed(4), asset, executionPrice)
}
