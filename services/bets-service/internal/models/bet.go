package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type BetDirection string
type BetStatus string

const (
	DirectionUp   BetDirection = "UP"
	DirectionDown BetDirection = "DOWN"

	StatusOpen      BetStatus = "OPEN"
	StatusWon       BetStatus = "WON"
	StatusLost      BetStatus = "LOST"
	StatusCancelled BetStatus = "CANCELLED"
)

type Bet struct {
	ID            int64           `json:"id"`
	UserID        int64           `json:"user_id"`
	Symbol        string          `json:"symbol"`
	Direction     BetDirection    `json:"direction"`
	StakeAmount   decimal.Decimal `json:"stake_amount"`
	OpenedPrice   decimal.Decimal `json:"opened_price"`
	ResolvedPrice decimal.Decimal `json:"resolved_price"` // 0 if not resolved
	Status        BetStatus       `json:"status"`
	OpenedAt      time.Time       `json:"opened_at"`
	ExpiresAt     time.Time       `json:"expires_at"`
	ResolvedAt    *time.Time      `json:"resolved_at,omitempty"`
}
