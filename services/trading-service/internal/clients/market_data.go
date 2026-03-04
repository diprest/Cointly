package clients

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/shopspring/decimal"
)

type MarketDataHTTP struct {
	BaseURL string
	Client  *http.Client
}

func NewMarketDataHTTP(url string) *MarketDataHTTP {
	return &MarketDataHTTP{
		BaseURL: url,
		Client:  &http.Client{Timeout: 5 * time.Second},
	}
}

type TickerResponse struct {
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price"`
}

func (m *MarketDataHTTP) GetPrice(symbol string) (decimal.Decimal, error) {
	url := fmt.Sprintf("%s/api/v1/market/ticker?symbol=%s", m.BaseURL, symbol)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return decimal.Zero, err
	}

	resp, err := m.Client.Do(req)
	if err != nil {
		return decimal.Zero, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return decimal.Zero, fmt.Errorf("market data error: %s", resp.Status)
	}

	var ticker TickerResponse
	if err := json.NewDecoder(resp.Body).Decode(&ticker); err != nil {
		return decimal.Zero, err
	}

	return decimal.NewFromFloat(ticker.Price), nil
}
