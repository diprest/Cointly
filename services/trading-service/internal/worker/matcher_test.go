package worker

import (
	"context"
	"testing"
	"time"
	"trading-service/internal/storage"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

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

func TestMatcher_PollAndMatch(t *testing.T) {
	db, mockDB, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	store := &storage.Storage{DB: db}
	pf := new(MockPortfolio)
	md := new(MockMarketData)
	matcher := NewMatcher(store, pf, md, nil)

	rows := sqlmock.NewRows([]string{"id", "user_id", "symbol", "side", "type", "price", "amount", "quote_amount", "status"}).
		AddRow(1, 101, "BTCUSDT", "BUY", "LIMIT", 50000.0, 0.1, 0.0, "NEW")
	mockDB.ExpectQuery("SELECT id, user_id, symbol, side, type, price, amount, quote_amount, status FROM orders WHERE status = 'NEW'").
		WillReturnRows(rows)
	md.On("GetPrice", "BTCUSDT").Return(decimal.NewFromInt(50000), nil)

	buyRows := sqlmock.NewRows([]string{"id", "user_id", "symbol", "side", "type", "price", "amount", "quote_amount", "status"}).
		AddRow(1, 101, "BTCUSDT", "BUY", "LIMIT", 50000.0, 0.1, 0.0, "NEW")
	mockDB.ExpectQuery("SELECT id, user_id, symbol, side, type, price, amount, quote_amount, status FROM orders WHERE symbol=\\$1 AND status='NEW' AND side='BUY'").
		WithArgs("BTCUSDT", decimal.NewFromInt(50000)).
		WillReturnRows(buyRows)
	pf.On("TransferFunds", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mockDB.ExpectExec("UPDATE orders SET status = 'FILLED', updated_at = \\$1 WHERE id = \\$2").
		WithArgs(sqlmock.AnyArg(), 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mockDB.ExpectQuery("SELECT id, user_id, symbol, side, type, price, amount, quote_amount, status FROM orders WHERE symbol=\\$1 AND status='NEW' AND side='SELL'").
		WithArgs("BTCUSDT", decimal.NewFromInt(50000)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	matcher.pollAndMatch()

	assert.NoError(t, mockDB.ExpectationsWereMet())
	pf.AssertExpectations(t)
	md.AssertExpectations(t)
}

func TestMatcher_StartPolling(t *testing.T) {
	db, mockDB, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	store := &storage.Storage{DB: db}
	pf := new(MockPortfolio)
	md := new(MockMarketData)
	matcher := NewMatcher(store, pf, md, nil)

	mockDB.ExpectQuery("SELECT id, user_id, symbol, side, type, price, amount, quote_amount, status FROM orders WHERE status = 'NEW'").
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(10 * time.Millisecond)
		cancel()
	}()

	matcher.StartPolling(ctx, 1*time.Millisecond)

	assert.NoError(t, mockDB.ExpectationsWereMet())
}
