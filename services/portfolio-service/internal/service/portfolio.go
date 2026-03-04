package service

import (
	"context"
	"errors"
	"log/slog"
	"portfolio-service/internal/models"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

var (
	ErrInsufficientFunds       = errors.New("insufficient available funds")
	ErrInsufficientLockedFunds = errors.New("insufficient locked funds to unlock")
)

type BalanceRepository interface {
	GetBalanceByAsset(ctx context.Context, userID int, asset string) (*models.Balance, error)
	GetBalancesByUserID(ctx context.Context, userID int) ([]models.Balance, error)
	AtomicallyLockFunds(ctx context.Context, userID int, asset string, amount decimal.Decimal) error
	AtomicallyUnlockFunds(ctx context.Context, userID int, asset string, amount decimal.Decimal) error
	AtomicallyTransferFunds(ctx context.Context, userID int, asset string, amount decimal.Decimal, cost decimal.Decimal, transferType string) error
	CreateBalance(ctx context.Context, balance *models.Balance) error
	AtomicallyChangeBalance(ctx context.Context, userID int, asset string, amount decimal.Decimal) error
	ResetBalance(ctx context.Context, userID int) error
}

type PortfolioService struct {
	portfolioRepo BalanceRepository
}

func NewPortfolioService(b BalanceRepository) *PortfolioService {
	return &PortfolioService{portfolioRepo: b}
}

func (s *PortfolioService) GetBalance(ctx context.Context, userID int, asset string) (*models.Balance, error) {
	slog.Info("GetBalance called", "userID", userID, "asset", asset)
	balance, err := s.portfolioRepo.GetBalanceByAsset(ctx, userID, asset)
	if err != nil {
		slog.Info("GetBalance error", "error", err)
		if errors.Is(err, gorm.ErrRecordNotFound) && asset == "USDT" {
			slog.Info("Creating initial balance for user", "userID", userID)
			initialBal := &models.Balance{
				UserID:    userID,
				Asset:     asset,
				Amount:    decimal.NewFromInt(10000),
				LockedBal: decimal.Zero,
			}
			if createErr := s.portfolioRepo.CreateBalance(ctx, initialBal); createErr != nil {
				slog.Error("Failed to create initial balance", "error", createErr)
				return nil, createErr
			}
			return initialBal, nil
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &models.Balance{UserID: userID, Asset: asset, Amount: decimal.Zero, LockedBal: decimal.Zero}, nil
		}
		return nil, err
	}
	return balance, nil
}

func (s *PortfolioService) GetPortfolio(ctx context.Context, userID int) ([]models.Balance, error) {
	return s.portfolioRepo.GetBalancesByUserID(ctx, userID)
}

func (s *PortfolioService) LockFunds(ctx context.Context, userID int, asset string, amount decimal.Decimal) error {
	if amount.IsNegative() {
		return errors.New("amount to lock must be positive")
	}
	return s.portfolioRepo.AtomicallyLockFunds(ctx, userID, asset, amount)
}

func (s *PortfolioService) UnlockFunds(ctx context.Context, userID int, asset string, amount decimal.Decimal) error {
	if amount.IsNegative() {
		return errors.New("amount to unlock must be positive")
	}
	return s.portfolioRepo.AtomicallyUnlockFunds(ctx, userID, asset, amount)
}

func (s *PortfolioService) TransferFunds(ctx context.Context, userID int, asset string, amount decimal.Decimal, cost decimal.Decimal, transferType string) error {
	if userID <= 0 || asset == "" {
		return errors.New("invalid transfer parameters")
	}
	return s.portfolioRepo.AtomicallyTransferFunds(ctx, userID, asset, amount, cost, transferType)
}

func (s *PortfolioService) ChangeBalance(ctx context.Context, userID int, asset string, amount decimal.Decimal) error {
	if userID <= 0 || asset == "" {
		return errors.New("invalid parameters")
	}
	return s.portfolioRepo.AtomicallyChangeBalance(ctx, userID, asset, amount)
}

func (s *PortfolioService) ResetBalance(ctx context.Context, userID int) error {
	return s.portfolioRepo.ResetBalance(ctx, userID)
}
