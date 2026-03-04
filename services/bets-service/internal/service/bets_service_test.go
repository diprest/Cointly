package service

import (
	"bets-service/internal/models"
	"context"
	"errors"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockBetRepo struct {
	mock.Mock
}

func (m *MockBetRepo) CreateBet(ctx context.Context, bet *models.Bet) error {
	args := m.Called(ctx, bet)
	return args.Error(0)
}

func (m *MockBetRepo) GetUserBets(ctx context.Context, userID int64) ([]models.Bet, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]models.Bet), args.Error(1)
}

func (m *MockBetRepo) GetExpiredOpenBets(ctx context.Context) ([]models.Bet, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.Bet), args.Error(1)
}

func (m *MockBetRepo) UpdateBetStatus(ctx context.Context, bet *models.Bet) error {
	args := m.Called(ctx, bet)
	return args.Error(0)
}

func (m *MockBetRepo) ResetBets(ctx context.Context, userID int64) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

type MockPortfolio struct {
	mock.Mock
}

func (m *MockPortfolio) LockFunds(ctx context.Context, userID int64, asset string, amount decimal.Decimal) error {
	args := m.Called(ctx, userID, asset, amount)
	return args.Error(0)
}

func (m *MockPortfolio) UnlockFunds(ctx context.Context, userID int64, asset string, amount decimal.Decimal) error {
	args := m.Called(ctx, userID, asset, amount)
	return args.Error(0)
}

func (m *MockPortfolio) ChangeBalance(ctx context.Context, userID int64, asset string, amount decimal.Decimal) error {
	args := m.Called(ctx, userID, asset, amount)
	return args.Error(0)
}

type MockMarket struct {
	mock.Mock
}

func (m *MockMarket) GetPrice(ctx context.Context, symbol string) (decimal.Decimal, error) {
	args := m.Called(ctx, symbol)
	return args.Get(0).(decimal.Decimal), args.Error(1)
}

func TestCreateBet_Success(t *testing.T) {
	repo := new(MockBetRepo)
	pf := new(MockPortfolio)
	md := new(MockMarket)
	svc := NewBetsService(repo, pf, md)

	ctx := context.Background()
	userID := int64(101)
	symbol := "BTCUSDT"
	stake := decimal.NewFromInt(100)
	price := decimal.NewFromInt(50000)

	md.On("GetPrice", ctx, symbol).Return(price, nil)
	pf.On("LockFunds", ctx, userID, "USDT", stake).Return(nil)
	repo.On("CreateBet", ctx, mock.AnythingOfType("*models.Bet")).Return(nil)

	bet, err := svc.CreateBet(ctx, userID, symbol, models.DirectionUp, stake, 60)

	assert.NoError(t, err)
	assert.NotNil(t, bet)
	assert.Equal(t, price, bet.OpenedPrice)
	assert.Equal(t, models.StatusOpen, bet.Status)

	md.AssertExpectations(t)
	pf.AssertExpectations(t)
	repo.AssertExpectations(t)
}

func TestCreateBet_InsufficientFunds(t *testing.T) {
	repo := new(MockBetRepo)
	pf := new(MockPortfolio)
	md := new(MockMarket)
	svc := NewBetsService(repo, pf, md)

	ctx := context.Background()
	userID := int64(101)
	symbol := "BTCUSDT"
	stake := decimal.NewFromInt(100)
	price := decimal.NewFromInt(50000)

	md.On("GetPrice", ctx, symbol).Return(price, nil)
	pf.On("LockFunds", ctx, userID, "USDT", stake).Return(errors.New("insufficient funds"))

	bet, err := svc.CreateBet(ctx, userID, symbol, models.DirectionUp, stake, 60)

	assert.Error(t, err)
	assert.Nil(t, bet)
	assert.Contains(t, err.Error(), "insufficient funds")

	repo.AssertNotCalled(t, "CreateBet", mock.Anything, mock.Anything)
}
