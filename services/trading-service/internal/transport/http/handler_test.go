package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"trading-service/internal/models"
	"trading-service/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) CreateOrder(order *models.Order) error {
	args := m.Called(order)
	return args.Error(0)
}

func (m *MockRepo) GetOrder(id uint) (*models.Order, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Order), args.Error(1)
}

func (m *MockRepo) UpdateOrderStatus(id uint, status string) error {
	args := m.Called(id, status)
	return args.Error(0)
}

func (m *MockRepo) GetUserOrders(userID int64) ([]models.Order, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.Order), args.Error(1)
}

func (m *MockRepo) GetAllOpenOrders() ([]models.Order, error) {
	args := m.Called()
	return args.Get(0).([]models.Order), args.Error(1)
}

func (m *MockRepo) GetActiveOrders() ([]models.Order, error) {
	args := m.Called()
	return args.Get(0).([]models.Order), args.Error(1)
}

func (m *MockRepo) ResetOrders(userID int64) error {
	args := m.Called(userID)
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

func (m *MockPortfolio) TransferFunds(ctx context.Context, userID int64, asset string, amount decimal.Decimal, cost decimal.Decimal, side string) error {
	args := m.Called(ctx, userID, asset, amount, cost, side)
	return args.Error(0)
}

type MockMarketData struct {
	mock.Mock
}

func (m *MockMarketData) GetPrice(symbol string) (decimal.Decimal, error) {
	args := m.Called(symbol)
	return args.Get(0).(decimal.Decimal), args.Error(1)
}

func TestHandler_CreateOrder(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := new(MockRepo)
	pf := new(MockPortfolio)
	md := new(MockMarketData)
	svc := service.NewTradingService(repo, pf, md)
	handler := NewHandler(svc)
	router := handler.InitRoutes()

	reqBody := models.Order{
		UserID: 101,
		Symbol: "BTCUSDT",
		Side:   "BUY",
		Type:   "LIMIT",
		Price:  decimal.NewFromInt(40000),
		Amount: decimal.NewFromFloat(0.1),
	}
	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/trading/orders", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	md.On("GetPrice", "BTCUSDT").Return(decimal.NewFromInt(50000), nil)
	repo.On("CreateOrder", mock.Anything).Return(nil)
	router.ServeHTTP(w, req)

	t.Logf("Response Code: %d", w.Code)
	t.Logf("Response Body: %s", w.Body.String())

	assert.Equal(t, http.StatusCreated, w.Code)

	pf.AssertExpectations(t)
	repo.AssertExpectations(t)
}

func TestHandler_GetOrders(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := new(MockRepo)
	pf := new(MockPortfolio)
	md := new(MockMarketData)
	svc := service.NewTradingService(repo, pf, md)
	handler := NewHandler(svc)
	router := handler.InitRoutes()

	req, _ := http.NewRequest("GET", "/api/v1/trading/orders?user_id=101", nil)
	w := httptest.NewRecorder()

	repo.On("GetUserOrders", int64(101)).Return([]models.Order{}, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	repo.AssertExpectations(t)
}

func TestHandler_CancelOrder(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := new(MockRepo)
	pf := new(MockPortfolio)
	md := new(MockMarketData)
	svc := service.NewTradingService(repo, pf, md)
	handler := NewHandler(svc)
	router := handler.InitRoutes()

	req, _ := http.NewRequest("DELETE", "/api/v1/trading/orders/1?user_id=101", nil)
	w := httptest.NewRecorder()

	order := &models.Order{ID: 1, UserID: 101, Status: "NEW", Symbol: "BTCUSDT", Side: "BUY", Price: decimal.NewFromInt(50000), Amount: decimal.NewFromFloat(0.1)}
	repo.On("GetOrder", uint(1)).Return(order, nil)
	repo.On("UpdateOrderStatus", uint(1), "CANCELLED").Return(nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	repo.AssertExpectations(t)
}
