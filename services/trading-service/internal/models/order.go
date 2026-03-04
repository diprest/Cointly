package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type Order struct {
	ID          uint            `json:"id" db:"id"`
	UserID      int64           `json:"user_id" db:"user_id"`
	Symbol      string          `json:"symbol" db:"symbol"`
	Side        string          `json:"side" db:"side"`
	Type        string          `json:"type" db:"type"`
	Price       decimal.Decimal `json:"price" db:"price"`
	Amount      decimal.Decimal `json:"amount" db:"amount"`
	QuoteAmount decimal.Decimal `json:"quote_amount" db:"quote_amount"`
	Status      string          `json:"status" db:"status"`
	CreatedAt   time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at" db:"updated_at"`
}

type PriceUpdate struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}
