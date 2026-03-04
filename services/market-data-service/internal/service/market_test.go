package service

import (
	"context"
	"errors"
	"market-data-service/internal/models"
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

func (m *MockPriceRepo) GetPrice(ctx context.Context, symbol string) (float64, error) {
	args := m.Called(ctx, symbol)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockPriceRepo) SetPrice(ctx context.Context, symbol string, price float64) error {
	args := m.Called(ctx, symbol, price)
	return args.Error(0)
}

func (m *MockPriceRepo) GetOldestPrice(ctx context.Context, symbol string) (float64, error) {
	args := m.Called(ctx, symbol)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockPriceRepo) CachePnL(ctx context.Context, symbol string, pnl float64, ttl time.Duration) error {
	args := m.Called(ctx, symbol, pnl, ttl)
	return args.Error(0)
}

func (m *MockPriceRepo) GetCachedPnL(ctx context.Context, symbol string) (float64, error) {
	args := m.Called(ctx, symbol)
	return args.Get(0).(float64), args.Error(1)
}

func TestGetTicker_Success(t *testing.T) {
	mockPriceRepo := new(MockPriceRepo)
	mockSymbolRepo := new(MockSymbolRepo)
	svc := NewMarketService(mockSymbolRepo, mockPriceRepo)

	mockPriceRepo.On("GetPrice", mock.Anything, "BTCUSDT").Return(50000.0, nil)

	result, err := svc.GetTicker(context.Background(), "BTCUSDT")

	assert.NoError(t, err)
	assert.Equal(t, 50000.0, result.Price)
	assert.Equal(t, "BTCUSDT", result.Symbol)

	mockPriceRepo.AssertExpectations(t)
}

func TestGetTicker_NotFound(t *testing.T) {
	mockPriceRepo := new(MockPriceRepo)
	svc := NewMarketService(nil, mockPriceRepo)

	mockPriceRepo.On("GetPrice", mock.Anything, "UNKNOWN").Return(0.0, errors.New("redis nil"))

	_, err := svc.GetTicker(context.Background(), "UNKNOWN")

	assert.Error(t, err)
}

func TestGetAllSymbols(t *testing.T) {
	mockSymbolRepo := new(MockSymbolRepo)
	mockPriceRepo := new(MockPriceRepo)
	svc := NewMarketService(mockSymbolRepo, mockPriceRepo)

	expectedSymbols := []models.Symbol{
		{Symbol: "BTCUSDT", BaseAsset: "BTC", Name: "Bitcoin"},
		{Symbol: "ETHUSDT", BaseAsset: "ETH", Name: "Ethereum"},
	}

	mockSymbolRepo.On("GetActiveSymbols").Return(expectedSymbols, nil)

	mockPriceRepo.On("GetPrice", mock.Anything, "BTCUSDT").Return(50000.0, nil)
	mockPriceRepo.On("GetCachedPnL", mock.Anything, "BTCUSDT").Return(0.0, errors.New("miss"))
	mockPriceRepo.On("GetOldestPrice", mock.Anything, "BTCUSDT").Return(0.0, errors.New("no history"))

	mockPriceRepo.On("GetPrice", mock.Anything, "ETHUSDT").Return(3000.0, nil)
	mockPriceRepo.On("GetCachedPnL", mock.Anything, "ETHUSDT").Return(0.0, errors.New("miss"))
	mockPriceRepo.On("GetOldestPrice", mock.Anything, "ETHUSDT").Return(0.0, errors.New("no history"))

	symbols, err := svc.GetAllSymbols(context.Background())

	assert.NoError(t, err)
	assert.Len(t, symbols, 2)
	assert.Equal(t, "BTCUSDT", symbols[0].Symbol)
	assert.Equal(t, 50000.0, symbols[0].Price)
	assert.Equal(t, 0.0, symbols[0].PnL)
}
