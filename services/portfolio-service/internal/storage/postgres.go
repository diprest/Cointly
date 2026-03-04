package storage

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"portfolio-service/internal/models"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PostgresDB struct {
	DB *gorm.DB
}

func NewPostgresDB(db *gorm.DB) *PostgresDB {
	return &PostgresDB{DB: db}
}

func (r *PostgresDB) GetBalanceByAsset(ctx context.Context, userID int, asset string) (*models.Balance, error) {
	var balance models.Balance
	err := r.DB.WithContext(ctx).
		Where("user_id = ? AND asset = ?", userID, asset).
		First(&balance).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &balance, nil
}

func (r *PostgresDB) GetBalancesByUserID(ctx context.Context, userID int) ([]models.Balance, error) {
	var balances []models.Balance
	err := r.DB.WithContext(ctx).Where("user_id = ?", userID).Find(&balances).Error
	return balances, err
}

func (r *PostgresDB) CreateBalance(ctx context.Context, balance *models.Balance) error {
	return r.DB.WithContext(ctx).Create(balance).Error
}

func (r *PostgresDB) AtomicallyLockFunds(ctx context.Context, userID int, asset string, amount decimal.Decimal) error {
	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var balance models.Balance
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("user_id = ? AND asset = ?", userID, asset).
			First(&balance).Error

		if err != nil {
			return fmt.Errorf("balance not found: %w", err)
		}

		available := balance.Amount.Sub(balance.LockedBal)
		if available.LessThan(amount) {
			return fmt.Errorf("insufficient funds")
		}

		balance.LockedBal = balance.LockedBal.Add(amount)
		return tx.Save(&balance).Error
	})
}

func (r *PostgresDB) AtomicallyUnlockFunds(ctx context.Context, userID int, asset string, amount decimal.Decimal) error {
	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var balance models.Balance
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("user_id = ? AND asset = ?", userID, asset).
			First(&balance).Error

		if err != nil {
			return err
		}

		if balance.LockedBal.LessThan(amount) {
			balance.LockedBal = decimal.Zero
		} else {
			balance.LockedBal = balance.LockedBal.Sub(amount)
		}
		return tx.Save(&balance).Error
	})
}

func (r *PostgresDB) AtomicallyTransferFunds(ctx context.Context, userID int, asset string, amount decimal.Decimal, cost decimal.Decimal, transferType string) error {
	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		quoteAsset := "USDT"

		var assetBal models.Balance
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("user_id = ? AND asset = ?", userID, asset).
			First(&assetBal).Error

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				assetBal = models.Balance{UserID: userID, Asset: asset, Amount: decimal.Zero, LockedBal: decimal.Zero}
				if err := tx.Create(&assetBal).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		}

		var quoteBal models.Balance
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("user_id = ? AND asset = ?", userID, quoteAsset).
			First(&quoteBal).Error; err != nil {
			return fmt.Errorf("failed to get quote balance: %w", err)
		}

		slog.Info("TRANSFER EXEC", "Side", transferType, "Amount(BTC)", amount, "Cost(USDT)", cost)

		switch transferType {
		case "BUY":
			quoteBal.Amount = quoteBal.Amount.Sub(cost)

			if quoteBal.Amount.LessThan(quoteBal.LockedBal) {
				return fmt.Errorf("insufficient available funds (locked by other operations)")
			}

			assetBal.Amount = assetBal.Amount.Add(amount)
			assetBal.TotalCost = assetBal.TotalCost.Add(cost)

		case "SELL":
			absAmount := amount.Abs()

			if assetBal.Amount.LessThanOrEqual(absAmount) {
				assetBal.TotalCost = decimal.Zero
			} else {
				remainingAmount := assetBal.Amount.Sub(absAmount)
				ratio := remainingAmount.Div(assetBal.Amount)
				assetBal.TotalCost = assetBal.TotalCost.Mul(ratio)
			}

			assetBal.Amount = assetBal.Amount.Sub(absAmount)

			if assetBal.Amount.LessThan(assetBal.LockedBal) {
				return fmt.Errorf("insufficient available asset funds (locked by other operations)")
			}

			quoteBal.Amount = quoteBal.Amount.Add(cost)

		default:
			return fmt.Errorf("unknown transfer type: %s", transferType)
		}

		if err := tx.Save(&assetBal).Error; err != nil {
			return err
		}
		if err := tx.Save(&quoteBal).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *PostgresDB) AtomicallyChangeBalance(ctx context.Context, userID int, asset string, amount decimal.Decimal) error {
	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var balance models.Balance
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("user_id = ? AND asset = ?", userID, asset).
			First(&balance).Error

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				if amount.IsPositive() {
					balance = models.Balance{
						UserID:    userID,
						Asset:     asset,
						Amount:    amount,
						LockedBal: decimal.Zero,
					}
					return tx.Create(&balance).Error
				}
				return fmt.Errorf("balance not found and cannot withdraw")
			}
			return err
		}

		if amount.IsNegative() {
			if balance.Amount.Add(amount).IsNegative() {
				return fmt.Errorf("insufficient funds for withdrawal")
			}
		}

		balance.Amount = balance.Amount.Add(amount)
		return tx.Save(&balance).Error
	})
}

func (r *PostgresDB) ResetBalance(ctx context.Context, userID int) error {
	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ?", userID).Delete(&models.Balance{}).Error; err != nil {
			return err
		}
		initialBal := models.Balance{
			UserID:    userID,
			Asset:     "USDT",
			Amount:    decimal.NewFromInt(10000),
			LockedBal: decimal.Zero,
			TotalCost: decimal.Zero,
		}
		if err := tx.Create(&initialBal).Error; err != nil {
			return err
		}

		return nil
	})
}
