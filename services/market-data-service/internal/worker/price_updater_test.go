package worker

import (
	"context"
	"encoding/json"
	"market-data-service/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
)

type MockSymbolProvider struct {
	mock.Mock
}

func (m *MockSymbolProvider) GetActiveSymbols() ([]models.Symbol, error) {
	args := m.Called()
	return args.Get(0).([]models.Symbol), args.Error(1)
}

type MockPriceSaver struct {
	mock.Mock
}

func (m *MockPriceSaver) SetPrice(ctx context.Context, symbol string, price float64) error {
	args := m.Called(ctx, symbol, price)
	return args.Error(0)
}

func (m *MockPriceSaver) AddPriceHistory(ctx context.Context, symbol string, price float64, timestamp int64) error {
	args := m.Called(ctx, symbol, price, timestamp)
	return args.Error(0)
}

func (m *MockPriceSaver) TrimHistory(ctx context.Context, symbol string, retention int64) error {
	args := m.Called(ctx, symbol, retention)
	return args.Error(0)
}

type MockPublisher struct {
	mock.Mock
}

func (m *MockPublisher) Publish(symbol string, price float64) error {
	args := m.Called(symbol, price)
	return args.Error(0)
}

func TestUpdatePrices_Success(t *testing.T) {
	mockResponse := []BinanceTicker{
		{Symbol: "BTCUSDT", Price: "50000.00"},
		{Symbol: "ETHUSDT", Price: "3000.00"},
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	pg := new(MockSymbolProvider)
	redis := new(MockPriceSaver)
	broker := new(MockPublisher)

	pg.On("GetActiveSymbols").Return([]models.Symbol{
		{Symbol: "BTCUSDT"},
		{Symbol: "ETHUSDT"},
	}, nil)

	redis.On("SetPrice", mock.Anything, "BTCUSDT", 50000.00).Return(nil)
	redis.On("SetPrice", mock.Anything, "ETHUSDT", 3000.00).Return(nil)

	redis.On("AddPriceHistory", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()
	redis.On("TrimHistory", mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()

	broker.On("Publish", "BTCUSDT", 50000.00).Return(nil)
	broker.On("Publish", "ETHUSDT", 3000.00).Return(nil)

	updater := NewPriceUpdater(pg, redis, broker, 1*time.Second)
	updater.SetAPIURL(server.URL)

	updater.updatePrices(context.Background())
	pg.AssertExpectations(t)
	redis.AssertExpectations(t)
	broker.AssertExpectations(t)
}

func TestUpdatePrices_EmptySymbols(t *testing.T) {
	pg := new(MockSymbolProvider)
	redis := new(MockPriceSaver)
	broker := new(MockPublisher)

	pg.On("GetActiveSymbols").Return([]models.Symbol{}, nil)

	updater := NewPriceUpdater(pg, redis, broker, 1*time.Second)
	updater.updatePrices(context.Background())

	pg.AssertExpectations(t)
	redis.AssertNotCalled(t, "SetPrice", mock.Anything, mock.Anything, mock.Anything)
}
