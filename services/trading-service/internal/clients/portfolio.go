package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/shopspring/decimal"
)

type PortfolioClient interface {
	LockFunds(ctx context.Context, userID int64, asset string, amount decimal.Decimal) error
	UnlockFunds(ctx context.Context, userID int64, asset string, amount decimal.Decimal) error
	TransferFunds(ctx context.Context, userID int64, asset string, amount decimal.Decimal, cost decimal.Decimal, side string) error
}

type MarketDataClient interface {
	GetPrice(symbol string) (decimal.Decimal, error)
}

type PortfolioHTTP struct {
	BaseURL string
	Client  *http.Client
}

func NewPortfolioHTTP(url string) *PortfolioHTTP {
	return &PortfolioHTTP{
		BaseURL: url,
		Client:  &http.Client{},
	}
}

func (p *PortfolioHTTP) LockFunds(ctx context.Context, userID int64, asset string, amount decimal.Decimal) error {
	payload := map[string]interface{}{
		"user_id": userID,
		"asset":   asset,
		"amount":  amount,
	}
	return p.sendRequest(ctx, "POST", "/lock", payload)
}

func (p *PortfolioHTTP) UnlockFunds(ctx context.Context, userID int64, asset string, amount decimal.Decimal) error {
	payload := map[string]interface{}{
		"user_id": userID,
		"asset":   asset,
		"amount":  amount,
	}
	return p.sendRequest(ctx, "POST", "/unlock", payload)
}

func (p *PortfolioHTTP) TransferFunds(ctx context.Context, userID int64, asset string, amount decimal.Decimal, cost decimal.Decimal, side string) error {
	payload := map[string]interface{}{
		"user_id": userID,
		"asset":   asset,
		"amount":  amount,
		"cost":    cost,
		"side":    side,
	}
	return p.sendRequest(ctx, "POST", "/transfer", payload)
}

func (p *PortfolioHTTP) sendRequest(ctx context.Context, method, endpoint string, payload interface{}) error {
	body, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, method, p.BaseURL+"/api/v1/portfolio"+endpoint, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return errors.New("portfolio service error status: " + resp.Status)
	}
	return nil
}
