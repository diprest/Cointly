package storage

import (
	"context"
	"log"
	"os"
	"testing"

	"portfolio-service/internal/models"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var testDB *gorm.DB
var testRepo *PostgresDB

func TestMain(m *testing.M) {
	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		dbURL = "host=localhost user=user password=password dbname=portfolio_db port=5432 sslmode=disable"
	}

	var err error
	testDB, err = gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Printf("Could not connect to test database: %v. Skipping tests.", err)
		os.Exit(0)
	}

	err = testDB.AutoMigrate(&models.Balance{})
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	testRepo = NewPostgresDB(testDB)

	code := m.Run()
	os.Exit(code)
}

func setupTestDB(t *testing.T) *gorm.DB {
	tx := testDB.Begin()
	t.Cleanup(func() {
		tx.Rollback()
	})
	return tx
}

func TestAtomicallyTransferFunds_MarketBuy_Success(t *testing.T) {
	txDB := setupTestDB(t)
	repo := NewPostgresDB(txDB)
	ctx := context.Background()
	userID := 900

	startUSDT := decimal.NewFromInt(1000)
	startBTC := decimal.NewFromFloat(0.5)

	txDB.Create(&models.Balance{
		UserID: userID, Asset: "USDT", Amount: startUSDT, LockedBal: decimal.Zero, TotalCost: decimal.Zero,
	})
	txDB.Create(&models.Balance{
		UserID: userID, Asset: "BTC", Amount: startBTC, LockedBal: decimal.Zero, TotalCost: decimal.Zero,
	})

	amountBTC := decimal.NewFromFloat(0.02)
	costUSDT := decimal.NewFromInt(500)

	txDB.Model(&models.Balance{}).
		Where("user_id = ? AND asset = ?", userID, "USDT").
		Update("locked_bal", costUSDT)

	err := repo.AtomicallyTransferFunds(ctx, userID, "BTC", amountBTC, costUSDT, "BUY")

	assert.NoError(t, err)

	usdtBal, _ := repo.GetBalanceByAsset(ctx, userID, "USDT")
	btcBal, _ := repo.GetBalanceByAsset(ctx, userID, "BTC")

	assert.Equal(t, decimal.NewFromInt(500).String(), usdtBal.Amount.String())
	assert.Equal(t, decimal.Zero.String(), usdtBal.LockedBal.String())

	assert.Equal(t, decimal.NewFromFloat(0.52).String(), btcBal.Amount.String())
	assert.Equal(t, decimal.Zero.String(), btcBal.LockedBal.String())
}

func TestAtomicallyTransferFunds_MarketSell_Success(t *testing.T) {
	txDB := setupTestDB(t)
	repo := NewPostgresDB(txDB)
	ctx := context.Background()
	userID := 901

	startUSDT := decimal.NewFromInt(100)
	startBTC := decimal.NewFromFloat(1.0)

	txDB.Create(&models.Balance{
		UserID: userID, Asset: "USDT", Amount: startUSDT, LockedBal: decimal.Zero, TotalCost: decimal.Zero,
	})
	txDB.Create(&models.Balance{
		UserID: userID, Asset: "BTC", Amount: startBTC, LockedBal: decimal.Zero, TotalCost: decimal.Zero,
	})

	amountBTC := decimal.NewFromFloat(0.5)
	costUSDT := decimal.NewFromInt(30000)

	txDB.Model(&models.Balance{}).
		Where("user_id = ? AND asset = ?", userID, "BTC").
		Update("locked_bal", amountBTC)

	err := repo.AtomicallyTransferFunds(ctx, userID, "BTC", amountBTC, costUSDT, "SELL")

	assert.NoError(t, err)

	usdtBal, _ := repo.GetBalanceByAsset(ctx, userID, "USDT")
	btcBal, _ := repo.GetBalanceByAsset(ctx, userID, "BTC")

	assert.Equal(t, decimal.NewFromInt(30100).String(), usdtBal.Amount.String())
	assert.Equal(t, decimal.Zero.String(), usdtBal.LockedBal.String())

	assert.Equal(t, decimal.NewFromFloat(0.5).String(), btcBal.Amount.String())
	assert.Equal(t, decimal.Zero.String(), btcBal.LockedBal.String())
}
