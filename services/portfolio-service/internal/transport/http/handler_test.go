package http

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"portfolio-service/internal/models"
	"portfolio-service/internal/service"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockBalanceRepo struct {
	mock.Mock
}

func (m *MockBalanceRepo) GetBalanceByAsset(ctx context.Context, userID int, asset string) (*models.Balance, error) {
	args := m.Called(ctx, userID, asset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Balance), args.Error(1)
}

func (m *MockBalanceRepo) GetBalancesByUserID(ctx context.Context, userID int) ([]models.Balance, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]models.Balance), args.Error(1)
}

func (m *MockBalanceRepo) AtomicallyLockFunds(ctx context.Context, userID int, asset string, amount decimal.Decimal) error {
	args := m.Called(ctx, userID, asset, amount)
	return args.Error(0)
}

func (m *MockBalanceRepo) AtomicallyUnlockFunds(ctx context.Context, userID int, asset string, amount decimal.Decimal) error {
	args := m.Called(ctx, userID, asset, amount)
	return args.Error(0)
}

func (m *MockBalanceRepo) AtomicallyTransferFunds(ctx context.Context, userID int, asset string, amount decimal.Decimal, cost decimal.Decimal, transferType string) error {
	args := m.Called(ctx, userID, asset, amount, cost, transferType)
	return args.Error(0)
}

func (m *MockBalanceRepo) CreateBalance(ctx context.Context, balance *models.Balance) error {
	args := m.Called(ctx, balance)
	return args.Error(0)
}

func (m *MockBalanceRepo) AtomicallyChangeBalance(ctx context.Context, userID int, asset string, amount decimal.Decimal) error {
	args := m.Called(ctx, userID, asset, amount)
	return args.Error(0)
}

func (m *MockBalanceRepo) ResetBalance(ctx context.Context, userID int) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func TestHandler_GetBalance(t *testing.T) {
	repo := new(MockBalanceRepo)
	svc := service.NewPortfolioService(repo)
	handler := NewHandler(svc)

	req, _ := http.NewRequest("GET", "/api/v1/portfolio/balance?user_id=101&asset=BTC", nil)
	w := httptest.NewRecorder()

	repo.On("GetBalanceByAsset", mock.Anything, 101, "BTC").Return(&models.Balance{UserID: 101, Asset: "BTC", Amount: decimal.NewFromInt(1)}, nil)

	handler.handleGetBalance(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	repo.AssertExpectations(t)
}

func TestHandler_GetPortfolio(t *testing.T) {
	repo := new(MockBalanceRepo)
	svc := service.NewPortfolioService(repo)
	handler := NewHandler(svc)

	req, _ := http.NewRequest("GET", "/api/v1/portfolio?user_id=101", nil)
	w := httptest.NewRecorder()

	repo.On("GetBalancesByUserID", mock.Anything, 101).Return([]models.Balance{}, nil)

	handler.handleGetPortfolio(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	repo.AssertExpectations(t)
}

func TestHandler_LockFunds(t *testing.T) {
	repo := new(MockBalanceRepo)
	svc := service.NewPortfolioService(repo)
	handler := NewHandler(svc)

	payload := `{"user_id": 101, "asset": "USDT", "amount": "100"}`
	req, _ := http.NewRequest("POST", "/api/v1/portfolio/lock", bytes.NewBufferString(payload))
	w := httptest.NewRecorder()

	repo.On("AtomicallyLockFunds", mock.Anything, 101, "USDT", decimal.NewFromInt(100)).Return(nil)

	handler.handleLockFunds(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	repo.AssertExpectations(t)
}

func TestHandler_UnlockFunds(t *testing.T) {
	repo := new(MockBalanceRepo)
	svc := service.NewPortfolioService(repo)
	handler := NewHandler(svc)

	payload := `{"user_id": 101, "asset": "USDT", "amount": "50"}`
	req, _ := http.NewRequest("POST", "/api/v1/portfolio/unlock", bytes.NewBufferString(payload))
	w := httptest.NewRecorder()

	repo.On("AtomicallyUnlockFunds", mock.Anything, 101, "USDT", decimal.NewFromInt(50)).Return(nil)

	handler.handleUnlockFunds(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	repo.AssertExpectations(t)
}

func TestHandler_TransferFunds(t *testing.T) {
	repo := new(MockBalanceRepo)
	svc := service.NewPortfolioService(repo)
	handler := NewHandler(svc)

	payload := `{"user_id": 101, "asset": "BTC", "amount": "0.1", "cost": "5000", "side": "BUY"}`
	req, _ := http.NewRequest("POST", "/api/v1/portfolio/transfer", bytes.NewBufferString(payload))
	w := httptest.NewRecorder()

	repo.On("AtomicallyTransferFunds", mock.Anything, 101, "BTC", decimal.NewFromFloat(0.1), decimal.NewFromInt(5000), "BUY").Return(nil)

	handler.handleTransferFunds(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	repo.AssertExpectations(t)
}
