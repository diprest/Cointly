package storage

import (
	"context"
	"fmt"
	"portfolio-service/internal/models"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestPnLCalculation(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}
	db.AutoMigrate(&models.Balance{})

	repo := NewPostgresDB(db)
	ctx := context.Background()
	userID := 999

	initialUSDT := decimal.NewFromInt(200000)
	err = repo.AtomicallyChangeBalance(ctx, userID, "USDT", initialUSDT)
	assert.NoError(t, err)
	usdt, _ := repo.GetBalanceByAsset(ctx, userID, "USDT")
	assert.Equal(t, initialUSDT.String(), usdt.Amount.String())

	amount1 := decimal.NewFromInt(1)
	price1 := decimal.NewFromInt(50000)
	cost1 := amount1.Mul(price1)

	err = repo.AtomicallyTransferFunds(ctx, userID, "BTC", amount1, cost1, "BUY")
	assert.NoError(t, err)

	btc, _ := repo.GetBalanceByAsset(ctx, userID, "BTC")
	assert.Equal(t, "1", btc.Amount.String())
	assert.Equal(t, "50000", btc.TotalCost.String())
	amount2 := decimal.NewFromInt(1)
	price2 := decimal.NewFromInt(60000)
	cost2 := amount2.Mul(price2)

	err = repo.AtomicallyTransferFunds(ctx, userID, "BTC", amount2, cost2, "BUY")
	assert.NoError(t, err)

	btc, _ = repo.GetBalanceByAsset(ctx, userID, "BTC")
	assert.Equal(t, "2", btc.Amount.String())
	assert.Equal(t, "110000", btc.TotalCost.String())

	amountSell := decimal.NewFromInt(1)
	priceSell := decimal.NewFromInt(70000)
	costSell := amountSell.Mul(priceSell)

	err = repo.AtomicallyTransferFunds(ctx, userID, "BTC", amountSell, costSell, "SELL")
	assert.NoError(t, err)

	btc, _ = repo.GetBalanceByAsset(ctx, userID, "BTC")
	assert.Equal(t, "1", btc.Amount.String())
	assert.Equal(t, "55000", btc.TotalCost.String())
	usdt, _ = repo.GetBalanceByAsset(ctx, userID, "USDT")
	assert.Equal(t, "160000", usdt.Amount.String())

	fmt.Println("PnL Test Passed Successfully!")
	fmt.Printf("Final BTC Amount: %s, TotalCost: %s\n", btc.Amount, btc.TotalCost)
	fmt.Printf("Final USDT Balance: %s\n", usdt.Amount)
}
