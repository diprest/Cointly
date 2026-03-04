package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/shopspring/decimal"
)

type PortfolioClient struct {
	BaseURL string
	Client  *http.Client
}

func NewPortfolioClient(url string) *PortfolioClient {
	return &PortfolioClient{
		BaseURL: url,
		Client:  &http.Client{},
	}
}

func (p *PortfolioClient) LockFunds(ctx context.Context, userID int64, asset string, amount decimal.Decimal) error {
	payload := map[string]interface{}{
		"user_id": userID,
		"asset":   asset,
		"amount":  amount,
	}
	return p.sendRequest(ctx, "POST", "/lock", payload)
}

func (p *PortfolioClient) UnlockFunds(ctx context.Context, userID int64, asset string, amount decimal.Decimal) error {
	payload := map[string]interface{}{
		"user_id": userID,
		"asset":   asset,
		"amount":  amount,
	}
	return p.sendRequest(ctx, "POST", "/unlock", payload)
}

func (p *PortfolioClient) ChangeBalance(ctx context.Context, userID int64, asset string, amount decimal.Decimal) error {
	payload := map[string]interface{}{
		"user_id": userID,
		"asset":   asset,
		"amount":  amount,
	}
	return p.sendRequest(ctx, "POST", "/balance/change", payload)
}

func (p *PortfolioClient) sendRequest(ctx context.Context, method, endpoint string, payload interface{}) error {
	body, _ := json.Marshal(payload)
	url := fmt.Sprintf("%s/api/v1/portfolio%s", p.BaseURL, endpoint)
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(body))
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
		return errors.New("portfolio service error: " + resp.Status)
	}
	return nil
}
