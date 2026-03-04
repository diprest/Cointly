package http

import (
	"context"
	"encoding/json"
	"errors"
	"market-data-service/internal/models"
	"market-data-service/internal/service"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockSymbolRepo struct {
	mock.Mock
}

func (m *MockSymbolRepo) GetActiveSymbols() ([]models.Symbol, error) {
	args := m.Called()
	return args.Get(0).([]models.Symbol), args.Error(1)
}

type MockPriceRepo struct {
	mock.Mock
}

func (m *MockPriceRepo) SetPrice(ctx context.Context, symbol string, price float64) error {
	args := m.Called(ctx, symbol, price)
	return args.Error(0)
}

func (m *MockPriceRepo) GetPrice(ctx context.Context, symbol string) (float64, error) {
	args := m.Called(ctx, symbol)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockPriceRepo) GetOldestPrice(ctx context.Context, symbol string) (float64, error) {
	args := m.Called(ctx, symbol)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockPriceRepo) CachePnL(ctx context.Context, symbol string, pnl float64, duration time.Duration) error {
	args := m.Called(ctx, symbol, pnl, duration)
	return args.Error(0)
}

func (m *MockPriceRepo) GetCachedPnL(ctx context.Context, symbol string) (float64, error) {
	args := m.Called(ctx, symbol)
	return args.Get(0).(float64), args.Error(1)
}

func TestHandler_GetSymbols(t *testing.T) {
	mockSymbolRepo := new(MockSymbolRepo)
	mockPriceRepo := new(MockPriceRepo)
	svc := service.NewMarketService(mockSymbolRepo, mockPriceRepo)
	handler := NewHandler(svc)

	mockSymbolRepo.On("GetActiveSymbols").Return([]models.Symbol{
		{Symbol: "BTCUSDT", Name: "Bitcoin"},
	}, nil)
	mockPriceRepo.On("GetPrice", mock.Anything, "BTCUSDT").Return(50000.0, nil)
	mockPriceRepo.On("GetCachedPnL", mock.Anything, "BTCUSDT").Return(5.0, nil)

	req, _ := http.NewRequest("GET", "/api/v1/market/symbols", nil)
	w := httptest.NewRecorder()

	handler.handleGetSymbols(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp []models.CoinInfo
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Len(t, resp, 1)
	assert.Equal(t, "BTCUSDT", resp[0].Symbol)
	assert.Equal(t, 50000.0, resp[0].Price)
	assert.Equal(t, 5.0, resp[0].PnL)
}

func TestHandler_GetSymbols_Error(t *testing.T) {
	mockSymbolRepo := new(MockSymbolRepo)
	mockPriceRepo := new(MockPriceRepo)
	svc := service.NewMarketService(mockSymbolRepo, mockPriceRepo)
	handler := NewHandler(svc)
	mockSymbolRepo.On("GetActiveSymbols").Return([]models.Symbol{}, errors.New("db error"))

	req, _ := http.NewRequest("GET", "/api/v1/market/symbols", nil)
	w := httptest.NewRecorder()

	handler.handleGetSymbols(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestHandler_RegisterRoutes(t *testing.T) {
	mockSymbolRepo := new(MockSymbolRepo)
	mockPriceRepo := new(MockPriceRepo)
	svc := service.NewMarketService(mockSymbolRepo, mockPriceRepo)
	handler := NewHandler(svc)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	req, _ := http.NewRequest("GET", "/api/v1/market/symbols", nil)
	_, pattern := mux.Handler(req)
	assert.NotEmpty(t, pattern)
}

func TestHandler_GetTicker(t *testing.T) {
	mockSymbolRepo := new(MockSymbolRepo)
	mockPriceRepo := new(MockPriceRepo)
	svc := service.NewMarketService(mockSymbolRepo, mockPriceRepo)
	handler := NewHandler(svc)

	mockPriceRepo.On("GetPrice", mock.Anything, "BTCUSDT").Return(50000.0, nil)

	req, _ := http.NewRequest("GET", "/api/v1/market/ticker?symbol=BTCUSDT", nil)
	w := httptest.NewRecorder()
	handler.handleGetTicker(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp models.TickerPrice
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "BTCUSDT", resp.Symbol)
	assert.Equal(t, 50000.0, resp.Price)
}

func TestHandler_GetTicker_NotFound(t *testing.T) {
	mockSymbolRepo := new(MockSymbolRepo)
	mockPriceRepo := new(MockPriceRepo)
	svc := service.NewMarketService(mockSymbolRepo, mockPriceRepo)
	handler := NewHandler(svc)

	mockPriceRepo.On("GetPrice", mock.Anything, "UNKNOWN").Return(0.0, errors.New("not found"))

	req, _ := http.NewRequest("GET", "/api/v1/market/ticker?symbol=UNKNOWN", nil)
	w := httptest.NewRecorder()

	handler.handleGetTicker(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
