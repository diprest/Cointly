package worker

import (
	"bets-service/internal/models"
	"bets-service/internal/service"
	"context"
	"testing"
	"time"

	"github.com/shopspring/decimal"
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

func TestResolver_ProcessExpiredBets(t *testing.T) {
	repo := new(MockBetRepo)
	pf := new(MockPortfolio)
	md := new(MockMarket)
	svc := service.NewBetsService(repo, pf, md)
	resolver := NewResolver(svc)

	bet := models.Bet{
		ID:          1,
		UserID:      101,
		Symbol:      "BTCUSDT",
		Direction:   models.DirectionUp,
		OpenedPrice: decimal.NewFromInt(50000),
		StakeAmount: decimal.NewFromInt(100),
		Status:      models.StatusOpen,
		ExpiresAt:   time.Now().Add(-1 * time.Minute), // Expired
	}

	repo.On("GetExpiredOpenBets", mock.Anything).Return([]models.Bet{bet}, nil)
	md.On("GetPrice", mock.Anything, "BTCUSDT").Return(decimal.NewFromInt(51000), nil) // Won
	pf.On("UnlockFunds", mock.Anything, int64(101), "USDT", decimal.NewFromInt(100)).Return(nil)
	pf.On("ChangeBalance", mock.Anything, int64(101), "USDT", decimal.NewFromInt(100)).Return(nil)
	repo.On("UpdateBetStatus", mock.Anything, mock.Anything).Return(nil)

	resolver.processExpiredBets(context.Background())

	repo.AssertExpectations(t)
	md.AssertExpectations(t)
	pf.AssertExpectations(t)
}
