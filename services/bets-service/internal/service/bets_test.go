package service

import (
	"bets-service/internal/models"
	"context"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestBetsService_CreateBet(t *testing.T) {
	repo := new(MockBetRepo)
	pf := new(MockPortfolio)
	md := new(MockMarket)
	svc := NewBetsService(repo, pf, md)

	md.On("GetPrice", mock.Anything, "BTCUSDT").Return(decimal.NewFromInt(50000), nil)
	pf.On("LockFunds", mock.Anything, int64(101), "USDT", decimal.NewFromInt(100)).Return(nil)
	repo.On("CreateBet", mock.Anything, mock.Anything).Return(nil)
	bet, err := svc.CreateBet(context.Background(), 101, "BTCUSDT", models.DirectionUp, decimal.NewFromInt(100), 60)
	assert.NoError(t, err)
	assert.NotNil(t, bet)
	assert.Equal(t, decimal.NewFromInt(50000), bet.OpenedPrice)

	repo.AssertExpectations(t)
	pf.AssertExpectations(t)
	md.AssertExpectations(t)
}

func TestBetsService_ResolveBet_Won(t *testing.T) {
	repo := new(MockBetRepo)
	pf := new(MockPortfolio)
	md := new(MockMarket)
	svc := NewBetsService(repo, pf, md)

	bet := &models.Bet{
		ID:          1,
		UserID:      101,
		Symbol:      "BTCUSDT",
		Direction:   models.DirectionUp,
		OpenedPrice: decimal.NewFromInt(50000),
		StakeAmount: decimal.NewFromInt(100),
		Status:      models.StatusOpen,
	}

	md.On("GetPrice", mock.Anything, "BTCUSDT").Return(decimal.NewFromInt(51000), nil) // Won
	pf.On("UnlockFunds", mock.Anything, int64(101), "USDT", decimal.NewFromInt(100)).Return(nil)
	pf.On("ChangeBalance", mock.Anything, int64(101), "USDT", decimal.NewFromInt(100)).Return(nil)
	repo.On("UpdateBetStatus", mock.Anything, bet).Return(nil)
	err := svc.ResolveBet(context.Background(), bet)
	assert.NoError(t, err)
	assert.Equal(t, models.StatusWon, bet.Status)

	repo.AssertExpectations(t)
	pf.AssertExpectations(t)
	md.AssertExpectations(t)
}

func TestBetsService_ResolveBet_Lost(t *testing.T) {
	repo := new(MockBetRepo)
	pf := new(MockPortfolio)
	md := new(MockMarket)
	svc := NewBetsService(repo, pf, md)

	bet := &models.Bet{
		ID:          1,
		UserID:      101,
		Symbol:      "BTCUSDT",
		Direction:   models.DirectionUp,
		OpenedPrice: decimal.NewFromInt(50000),
		StakeAmount: decimal.NewFromInt(100),
		Status:      models.StatusOpen,
	}

	md.On("GetPrice", mock.Anything, "BTCUSDT").Return(decimal.NewFromInt(49000), nil) // Lost
	pf.On("UnlockFunds", mock.Anything, int64(101), "USDT", decimal.NewFromInt(100)).Return(nil)
	pf.On("ChangeBalance", mock.Anything, int64(101), "USDT", decimal.NewFromInt(-100)).Return(nil)
	repo.On("UpdateBetStatus", mock.Anything, bet).Return(nil)
	err := svc.ResolveBet(context.Background(), bet)
	assert.NoError(t, err)
	assert.Equal(t, models.StatusLost, bet.Status)

	repo.AssertExpectations(t)
	pf.AssertExpectations(t)
	md.AssertExpectations(t)
}
