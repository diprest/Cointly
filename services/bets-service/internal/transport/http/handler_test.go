package http

import (
	"bets-service/internal/models"
	"bets-service/internal/service"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
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

func TestHandler_CreateBet(t *testing.T) {
	repo := new(MockBetRepo)
	pf := new(MockPortfolio)
	md := new(MockMarket)
	svc := service.NewBetsService(repo, pf, md)
	handler := NewHandler(svc)

	reqBody := CreateBetRequest{
		UserID:      101,
		Symbol:      "BTCUSDT",
		Direction:   "UP",
		StakeAmount: decimal.NewFromInt(100),
		DurationSec: 60,
	}
	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/bets/", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	md.On("GetPrice", mock.Anything, "BTCUSDT").Return(decimal.NewFromInt(50000), nil)
	pf.On("LockFunds", mock.Anything, int64(101), "USDT", decimal.NewFromInt(100)).Return(nil)
	repo.On("CreateBet", mock.Anything, mock.Anything).Return(nil)

	handler.createBet(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	pf.AssertExpectations(t)
	repo.AssertExpectations(t)
}

func TestHandler_CreateBet_Error(t *testing.T) {
	repo := new(MockBetRepo)
	pf := new(MockPortfolio)
	md := new(MockMarket)
	svc := service.NewBetsService(repo, pf, md)
	handler := NewHandler(svc)

	reqBody := CreateBetRequest{
		UserID:      101,
		Symbol:      "BTCUSDT",
		Direction:   "UP",
		StakeAmount: decimal.NewFromInt(100),
		DurationSec: 60,
	}
	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/bets/", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	md.On("GetPrice", mock.Anything, "BTCUSDT").Return(decimal.NewFromInt(50000), nil)
	pf.On("LockFunds", mock.Anything, int64(101), "USDT", decimal.NewFromInt(100)).Return(errors.New("insufficient funds"))

	handler.createBet(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	pf.AssertExpectations(t)
}

func TestHandler_GetUserBets(t *testing.T) {
	repo := new(MockBetRepo)
	pf := new(MockPortfolio)
	md := new(MockMarket)
	svc := service.NewBetsService(repo, pf, md)
	handler := NewHandler(svc)

	req, _ := http.NewRequest("GET", "/api/v1/bets/?user_id=101", nil)
	w := httptest.NewRecorder()

	repo.On("GetUserBets", mock.Anything, int64(101)).Return([]models.Bet{}, nil)

	handler.getUserBets(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	repo.AssertExpectations(t)
}

func TestHandler_GetUserBets_Error(t *testing.T) {
	repo := new(MockBetRepo)
	pf := new(MockPortfolio)
	md := new(MockMarket)
	svc := service.NewBetsService(repo, pf, md)
	handler := NewHandler(svc)

	req, _ := http.NewRequest("GET", "/api/v1/bets/?user_id=101", nil)
	w := httptest.NewRecorder()

	repo.On("GetUserBets", mock.Anything, int64(101)).Return([]models.Bet{}, errors.New("db error"))

	handler.getUserBets(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	repo.AssertExpectations(t)
}

func TestHandler_InitRoutes(t *testing.T) {
	repo := new(MockBetRepo)
	pf := new(MockPortfolio)
	md := new(MockMarket)
	svc := service.NewBetsService(repo, pf, md)
	handler := NewHandler(svc)

	r := handler.InitRoutes()
	assert.NotNil(t, r)
}

func TestHandler_ResetBets(t *testing.T) {
	repo := new(MockBetRepo)
	pf := new(MockPortfolio)
	md := new(MockMarket)
	svc := service.NewBetsService(repo, pf, md)
	handler := NewHandler(svc)

	req, _ := http.NewRequest("POST", "/api/v1/bets/reset?user_id=101", nil)
	w := httptest.NewRecorder()

	repo.On("ResetBets", mock.Anything, int64(101)).Return(nil)

	handler.handleResetBets(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	repo.AssertExpectations(t)
}

func TestHandler_ResetBets_Error(t *testing.T) {
	repo := new(MockBetRepo)
	pf := new(MockPortfolio)
	md := new(MockMarket)
	svc := service.NewBetsService(repo, pf, md)
	handler := NewHandler(svc)

	req, _ := http.NewRequest("POST", "/api/v1/bets/reset?user_id=101", nil)
	w := httptest.NewRecorder()

	repo.On("ResetBets", mock.Anything, int64(101)).Return(errors.New("db error"))

	handler.handleResetBets(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	repo.AssertExpectations(t)
}
