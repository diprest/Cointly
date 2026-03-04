package clients

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/shopspring/decimal"
)

type MarketDataClient struct {
	BaseURL string
	Client  *http.Client
}

func NewMarketDataClient(url string) *MarketDataClient {
	return &MarketDataClient{
		BaseURL: url,
		Client:  &http.Client{},
	}
}

type PriceResponse struct {
	Symbol string          `json:"symbol"`
	Price  decimal.Decimal `json:"price"`
}

func (m *MarketDataClient) GetPrice(ctx context.Context, symbol string) (decimal.Decimal, error) {
	url := fmt.Sprintf("%s/api/v1/market/ticker?symbol=%s", m.BaseURL, symbol)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return decimal.Zero, err
	}

	resp, err := m.Client.Do(req)
	if err != nil {
		return decimal.Zero, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return decimal.Zero, errors.New("failed to get price")
	}

	var p PriceResponse
	if err := json.NewDecoder(resp.Body).Decode(&p); err != nil {
		return decimal.Zero, err
	}
	return p.Price, nil
}
