package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"bets-service/internal/models"

	"github.com/shopspring/decimal"
)

type BetRepository interface {
	CreateBet(ctx context.Context, bet *models.Bet) error
	GetUserBets(ctx context.Context, userID int64) ([]models.Bet, error)
	GetExpiredOpenBets(ctx context.Context) ([]models.Bet, error)
	UpdateBetStatus(ctx context.Context, bet *models.Bet) error
	ResetBets(ctx context.Context, userID int64) error
}

type PortfolioClient interface {
	LockFunds(ctx context.Context, userID int64, asset string, amount decimal.Decimal) error
	UnlockFunds(ctx context.Context, userID int64, asset string, amount decimal.Decimal) error
	ChangeBalance(ctx context.Context, userID int64, asset string, amount decimal.Decimal) error
}

type MarketDataClient interface {
	GetPrice(ctx context.Context, symbol string) (decimal.Decimal, error)
}

type BetsService struct {
	Repo      BetRepository
	Portfolio PortfolioClient
	Market    MarketDataClient
}

func NewBetsService(repo BetRepository, pf PortfolioClient, md MarketDataClient) *BetsService {
	return &BetsService{
		Repo:      repo,
		Portfolio: pf,
		Market:    md,
	}
}

func (s *BetsService) CreateBet(ctx context.Context, userID int64, symbol string, direction models.BetDirection, stake decimal.Decimal, durationSec int) (*models.Bet, error) {
	price, err := s.Market.GetPrice(ctx, symbol)
	if err != nil {
		return nil, fmt.Errorf("failed to get price: %w", err)
	}

	if err := s.Portfolio.LockFunds(ctx, userID, "USDT", stake); err != nil {
		return nil, fmt.Errorf("insufficient funds or lock failed: %w", err)
	}

	bet := &models.Bet{
		UserID:      userID,
		Symbol:      symbol,
		Direction:   direction,
		StakeAmount: stake,
		OpenedPrice: price,
		Status:      models.StatusOpen,
		OpenedAt:    time.Now(),
		ExpiresAt:   time.Now().Add(time.Duration(durationSec) * time.Second),
	}

	if err := s.Repo.CreateBet(ctx, bet); err != nil {
		_ = s.Portfolio.UnlockFunds(ctx, userID, "USDT", stake)
		return nil, err
	}

	return bet, nil
}

func (s *BetsService) GetUserBets(ctx context.Context, userID int64) ([]models.Bet, error) {
	return s.Repo.GetUserBets(ctx, userID)
}

func (s *BetsService) ResolveBet(ctx context.Context, bet *models.Bet) error {
	if bet.Status != models.StatusOpen {
		return nil
	}

	currentPrice, err := s.Market.GetPrice(ctx, bet.Symbol)
	if err != nil {
		return err
	}

	bet.ResolvedPrice = currentPrice
	now := time.Now()
	bet.ResolvedAt = &now

	won := false
	if bet.Direction == models.DirectionUp {
		if currentPrice.GreaterThan(bet.OpenedPrice) {
			won = true
		}
	} else if bet.Direction == models.DirectionDown {
		if currentPrice.LessThan(bet.OpenedPrice) {
			won = true
		}
	}

	if won {
		bet.Status = models.StatusWon
		if err := s.Portfolio.UnlockFunds(ctx, bet.UserID, "USDT", bet.StakeAmount); err != nil {
			slog.Error("Failed to unlock funds for won bet", "bet_id", bet.ID, "error", err)
			return err
		}
		if err := s.Portfolio.ChangeBalance(ctx, bet.UserID, "USDT", bet.StakeAmount); err != nil {
			slog.Error("Failed to add winnings", "bet_id", bet.ID, "error", err)
			return err
		}

	} else {
		bet.Status = models.StatusLost
		if err := s.Portfolio.UnlockFunds(ctx, bet.UserID, "USDT", bet.StakeAmount); err != nil {
			slog.Error("Failed to unlock funds for lost bet", "bet_id", bet.ID, "error", err)
			return err
		}
		if err := s.Portfolio.ChangeBalance(ctx, bet.UserID, "USDT", bet.StakeAmount.Neg()); err != nil {
			slog.Error("Failed to burn stake", "bet_id", bet.ID, "error", err)
			return err
		}
	}

	return s.Repo.UpdateBetStatus(ctx, bet)
}

func (s *BetsService) ResetBets(ctx context.Context, userID int64) error {
	return s.Repo.ResetBets(ctx, userID)
}
