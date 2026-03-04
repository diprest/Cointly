package models

import "github.com/shopspring/decimal"

type Balance struct {
	ID        int             `json:"id" gorm:"primaryKey"`
	UserID    int             `json:"user_id" gorm:"column:user_id"`
	Asset     string          `json:"asset" gorm:"column:asset"`
	Amount    decimal.Decimal `json:"amount" gorm:"column:amount;type:numeric"`
	LockedBal decimal.Decimal `json:"locked_bal" gorm:"column:locked_bal;type:numeric"`
	TotalCost decimal.Decimal `json:"total_cost" gorm:"column:total_cost;type:numeric"`
}

type LockUnlockRequest struct {
	UserID int    `json:"user_id"`
	Asset  string `json:"asset"`
	Amount string `json:"amount"`
}

type TransferRequest struct {
	UserID int    `json:"user_id"`
	Asset  string `json:"asset"`
	Amount string `json:"amount"`
	Cost   string `json:"cost"`
	Side   string `json:"side"`
}
