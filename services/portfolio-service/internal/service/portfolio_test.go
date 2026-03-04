package service

import (
	"context"
	"portfolio-service/internal/models"
	"testing"

	"gorm.io/gorm"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockBalanceRepository struct {
	mock.Mock
}

func (m *MockBalanceRepository) GetBalanceByAsset(ctx context.Context, userID int, asset string) (*models.Balance, error) {
	args := m.Called(ctx, userID, asset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Balance), args.Error(1)
}

func (m *MockBalanceRepository) GetBalancesByUserID(ctx context.Context, userID int) ([]models.Balance, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Balance), args.Error(1)
}

func (m *MockBalanceRepository) AtomicallyLockFunds(ctx context.Context, userID int, asset string, amount decimal.Decimal) error {
	args := m.Called(ctx, userID, asset, amount)
	return args.Error(0)
}

func (m *MockBalanceRepository) AtomicallyUnlockFunds(ctx context.Context, userID int, asset string, amount decimal.Decimal) error {
	args := m.Called(ctx, userID, asset, amount)
	return args.Error(0)
}

func (m *MockBalanceRepository) AtomicallyTransferFunds(ctx context.Context, userID int, asset string, amount decimal.Decimal, cost decimal.Decimal, transferType string) error {
	args := m.Called(ctx, userID, asset, amount, cost, transferType)
	return args.Error(0)
}

func (m *MockBalanceRepository) CreateBalance(ctx context.Context, balance *models.Balance) error {
	args := m.Called(ctx, balance)
	return args.Error(0)
}

func (m *MockBalanceRepository) AtomicallyChangeBalance(ctx context.Context, userID int, asset string, amount decimal.Decimal) error {
	args := m.Called(ctx, userID, asset, amount)
	return args.Error(0)
}

func (m *MockBalanceRepository) ResetBalance(ctx context.Context, userID int) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func TestGetBalance_NotFound_ReturnZeroBalance(t *testing.T) {
	repo := new(MockBalanceRepository)
	svc := NewPortfolioService(repo)
	ctx := context.Background()
	userID := 102
	asset := "BTC"

	repo.On("GetBalanceByAsset", ctx, userID, asset).Return(nil, gorm.ErrRecordNotFound).Once()

	balance, err := svc.GetBalance(ctx, userID, asset)

	assert.NoError(t, err, "Сервис должен обработать sql.ErrNoRows без возврата ошибки")
	assert.Equal(t, userID, balance.UserID)
	assert.Equal(t, asset, balance.Asset)
	assert.True(t, balance.Amount.IsZero(), "Amount должен быть 0")
	repo.AssertExpectations(t)
}

func TestLockFunds_NegativeAmount_Rejected(t *testing.T) {
	repo := new(MockBalanceRepository)
	svc := NewPortfolioService(repo)
	ctx := context.Background()
	userID := 103
	asset := "USDT"
	amount := decimal.NewFromInt(-100)

	repo.AssertNotCalled(t, "AtomicallyLockFunds", mock.Anything, mock.Anything, mock.Anything, mock.Anything)

	err := svc.LockFunds(ctx, userID, asset, amount)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "amount to lock must be positive")
}

func TestUnlockFunds_NegativeAmount_Rejected(t *testing.T) {
	repo := new(MockBalanceRepository)
	svc := NewPortfolioService(repo)
	ctx := context.Background()
	userID := 104
	asset := "USDT"
	amount := decimal.NewFromInt(-100)

	repo.AssertNotCalled(t, "AtomicallyUnlockFunds", mock.Anything, mock.Anything, mock.Anything, mock.Anything)

	err := svc.UnlockFunds(ctx, userID, asset, amount)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "amount to unlock must be positive")
}

func TestTransferFunds_Success(t *testing.T) {
	repo := new(MockBalanceRepository)
	svc := NewPortfolioService(repo)
	ctx := context.Background()
	userID := 105
	asset := "BTC"
	amount := decimal.NewFromFloat(0.1)
	cost := decimal.NewFromInt(5000)
	transferType := "BUY"

	repo.On("AtomicallyTransferFunds", ctx, userID, asset, amount, cost, transferType).Return(nil).Once()

	err := svc.TransferFunds(ctx, userID, asset, amount, cost, transferType)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestTransferFunds_InvalidInput(t *testing.T) {
	repo := new(MockBalanceRepository)
	svc := NewPortfolioService(repo)
	repo.AssertNotCalled(t, "AtomicallyTransferFunds", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)

	err := svc.TransferFunds(context.Background(), 0, "BTC", decimal.NewFromInt(1), decimal.NewFromInt(1), "BUY")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid transfer parameters")
	err = svc.TransferFunds(context.Background(), 105, "", decimal.NewFromInt(1), decimal.NewFromInt(1), "BUY")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid transfer parameters")
}
