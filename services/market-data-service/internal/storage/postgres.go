package storage

import (
	"log/slog"
	"market-data-service/internal/config"
	"market-data-service/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresDB struct {
	DB *gorm.DB
}

func NewPostgresDB(dsn string, initialSymbols []config.SymbolConfigItem) (*PostgresDB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&models.Symbol{}); err != nil {
		return nil, err
	}

	seedData(db, initialSymbols)
	return &PostgresDB{DB: db}, nil
}

func seedData(db *gorm.DB, initialSymbols []config.SymbolConfigItem) {
	slog.Info("Seeding/Updating symbols in DB")
	for _, s := range initialSymbols {
		symbol := models.Symbol{
			Symbol:     s.Symbol,
			Name:       s.Name,
			BaseAsset:  s.Symbol[:len(s.Symbol)-4],
			QuoteAsset: "USDT",
			IsActive:   true,
		}
		if err := db.Where(models.Symbol{Symbol: s.Symbol}).Assign(models.Symbol{Name: s.Name}).FirstOrCreate(&symbol).Error; err != nil {
			slog.Error("Failed to seed symbol", "symbol", s.Symbol, "error", err)
		}
	}
}

func (p *PostgresDB) GetActiveSymbols() ([]models.Symbol, error) {
	var symbols []models.Symbol
	result := p.DB.Where("is_active = ?", true).Find(&symbols)
	return symbols, result.Error
}
