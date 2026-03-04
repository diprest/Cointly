package models

import "time"

type Symbol struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Symbol     string    `gorm:"uniqueIndex;not null" json:"symbol"`
	Name       string    `json:"name"`
	BaseAsset  string    `json:"base_asset"`
	QuoteAsset string    `json:"quote_asset"`
	IsActive   bool      `json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
}

type TickerPrice struct {
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price"`
}

type CoinInfo struct {
	Symbol string  `json:"symbol"`
	Name   string  `json:"name"`
	Price  float64 `json:"price"`
	PnL    float64 `json:"pnl"`
}
