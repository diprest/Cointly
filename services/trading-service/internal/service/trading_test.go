package service

import (
	"context"
	"errors"
	"testing"
	"trading-service/internal/models"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepo struct{ mock.Mock }

func (m *MockRepo) CreateOrder(o *models.Order) error {
	args := m.Called(o)
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
	if args.Get(0) == nil {
		return []models.Order{}, args.Error(1)
	}
	return args.Get(0).([]models.Order), args.Error(1)
}

func (m *MockRepo) GetActiveOrders() ([]models.Order, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return []models.Order{}, args.Error(1)
	}
	return args.Get(0).([]models.Order), args.Error(1)
}

func (m *MockRepo) ResetOrders(userID int64) error {
	args := m.Called(userID)
	return args.Error(0)
}

type MockPortfolio struct{ mock.Mock }

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

type MockMarketData struct{ mock.Mock }

func (m *MockMarketData) GetPrice(symbol string) (decimal.Decimal, error) {
	args := m.Called(symbol)
	return args.Get(0).(decimal.Decimal), args.Error(1)
}

func TestCreateOrder_MarketBuyByQuote_ImmediateExecution(t *testing.T) {
	repo := new(MockRepo)
	pf := new(MockPortfolio)
	md := new(MockMarketData)
	svc := NewTradingService(repo, pf, md)

	quoteAmount := decimal.NewFromInt(1000)
	price := decimal.NewFromInt(50000)

	order := &models.Order{
		UserID:      101,
		Symbol:      "BTCUSDT",
		Side:        "BUY",
		Type:        "MARKET",
		QuoteAmount: quoteAmount,
	}

	md.On("GetPrice", "BTCUSDT").Return(price, nil)
	pf.On("TransferFunds", mock.Anything, int64(101), "BTC", mock.Anything, mock.Anything, "BUY").Return(nil)
	repo.On("CreateOrder", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		o := args.Get(0).(*models.Order)
		assert.Equal(t, "FILLED", o.Status)
		assert.Equal(t, price, o.Price)
	})

	err := svc.CreateOrder(order)
	assert.NoError(t, err)

	pf.AssertExpectations(t)
	repo.AssertExpectations(t)
	md.AssertExpectations(t)
}

func TestCreateOrder_MarketBuyByQuantity_ImmediateExecution(t *testing.T) {
	repo := new(MockRepo)
	pf := new(MockPortfolio)
	md := new(MockMarketData)
	svc := NewTradingService(repo, pf, md)

	amount := decimal.NewFromFloat(0.1)
	price := decimal.NewFromInt(60000)

	order := &models.Order{
		UserID: 102,
		Symbol: "BTCUSDT",
		Side:   "BUY",
		Type:   "MARKET",
		Amount: amount,
		Price:  price,
	}

	md.On("GetPrice", "BTCUSDT").Return(price, nil)
	pf.On("TransferFunds", mock.Anything, int64(102), "BTC", mock.Anything, mock.Anything, "BUY").Return(nil)
	repo.On("CreateOrder", mock.Anything).Return(nil)
	err := svc.CreateOrder(order)
	assert.NoError(t, err)

	pf.AssertExpectations(t)
	repo.AssertExpectations(t)
	md.AssertExpectations(t)
}

func TestCreateOrder_MarketSell_ImmediateExecution(t *testing.T) {
	repo := new(MockRepo)
	pf := new(MockPortfolio)
	md := new(MockMarketData)
	svc := NewTradingService(repo, pf, md)

	amount := decimal.NewFromFloat(0.5)
	price := decimal.NewFromInt(40000)

	order := &models.Order{
		UserID: 103,
		Symbol: "BTCUSDT",
		Side:   "SELL",
		Type:   "MARKET",
		Amount: amount,
	}

	md.On("GetPrice", "BTCUSDT").Return(price, nil)
	pf.On("TransferFunds", mock.Anything, int64(103), "BTC", mock.Anything, mock.Anything, "SELL").Return(nil)
	repo.On("CreateOrder", mock.Anything).Return(nil)
	err := svc.CreateOrder(order)
	assert.NoError(t, err)

	pf.AssertExpectations(t)
	repo.AssertExpectations(t)
	md.AssertExpectations(t)
}

func TestCreateOrder_TransferFundsFailure_OrderNotCreated(t *testing.T) {
	repo := new(MockRepo)
	pf := new(MockPortfolio)
	md := new(MockMarketData)
	svc := NewTradingService(repo, pf, md)

	quoteAmount := decimal.NewFromInt(10)
	price := decimal.NewFromInt(100)

	order := &models.Order{
		UserID: 104, Symbol: "BTCUSDT", Side: "BUY", Type: "MARKET",
		QuoteAmount: quoteAmount,
	}

	md.On("GetPrice", "BTCUSDT").Return(price, nil)
	pf.On("TransferFunds", mock.Anything, int64(104), "BTC", mock.Anything, quoteAmount, "BUY").Return(errors.New("insufficient funds")).Once()
	repo.AssertNotCalled(t, "CreateOrder", mock.Anything)

	err := svc.CreateOrder(order)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to execute order")

	pf.AssertExpectations(t)
	md.AssertExpectations(t)
}

func TestCancelOrder_UnlockQuoteAsset_Success(t *testing.T) {
	repo := new(MockRepo)
	pf := new(MockPortfolio)
	md := new(MockMarketData)
	svc := NewTradingService(repo, pf, md)

	orderID := uint(20)
	quoteAmount := decimal.NewFromInt(700)

	mockOrder := &models.Order{
		ID: orderID, UserID: 105, Symbol: "BTCUSDT", Side: "BUY",
		Type: "MARKET", Status: "NEW", QuoteAmount: quoteAmount,
	}

	repo.On("GetOrder", orderID).Return(mockOrder, nil).Once()
	repo.On("UpdateOrderStatus", orderID, "CANCELLED").Return(nil).Once()

	err := svc.CancelOrder(context.Background(), orderID, 105)
	assert.NoError(t, err)

	pf.AssertExpectations(t)
	repo.AssertExpectations(t)
}

func TestCancelOrder_UnlockQuantityAsset_Success(t *testing.T) {
	repo := new(MockRepo)
	pf := new(MockPortfolio)
	md := new(MockMarketData)
	svc := NewTradingService(repo, pf, md)

	orderID := uint(21)
	userID := int64(106)
	amount := decimal.NewFromFloat(0.5)

	mockOrder := &models.Order{
		ID: orderID, UserID: userID, Symbol: "BTCUSDT", Side: "SELL",
		Type: "LIMIT", Status: "NEW", Amount: amount,
	}

	repo.On("GetOrder", orderID).Return(mockOrder, nil).Once()
	repo.On("UpdateOrderStatus", orderID, "CANCELLED").Return(nil).Once()

	err := svc.CancelOrder(context.Background(), orderID, 106)
	assert.NoError(t, err)

	pf.AssertExpectations(t)
	repo.AssertExpectations(t)
}
