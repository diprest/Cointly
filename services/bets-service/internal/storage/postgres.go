package storage

import (
	"context"
	"database/sql"

	"bets-service/internal/models"

	_ "github.com/jackc/pgx/v5/stdlib" // pgx driver
)

type Storage struct {
	DB *sql.DB
}

func NewStorage(dsn string) (*Storage, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &Storage{DB: db}, nil
}

func (s *Storage) RunMigrations() error {
	query := `
	CREATE TABLE IF NOT EXISTS price_bets (
		id SERIAL PRIMARY KEY,
		user_id BIGINT NOT NULL,
		symbol VARCHAR(20) NOT NULL,
		direction VARCHAR(10) NOT NULL, -- 'UP' or 'DOWN'
		stake_amount DECIMAL(20, 8) NOT NULL,
		opened_price DECIMAL(20, 8) NOT NULL,
		resolved_price DECIMAL(20, 8) DEFAULT 0,
		status VARCHAR(20) NOT NULL, -- 'OPEN', 'WON', 'LOST', 'CANCELLED'
		opened_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
		resolved_at TIMESTAMP WITH TIME ZONE
	);
	`
	_, err := s.DB.Exec(query)
	return err
}

func (s *Storage) CreateBet(ctx context.Context, bet *models.Bet) error {
	query := `
		INSERT INTO price_bets (user_id, symbol, direction, stake_amount, opened_price, status, opened_at, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`
	return s.DB.QueryRowContext(ctx, query,
		bet.UserID, bet.Symbol, bet.Direction, bet.StakeAmount, bet.OpenedPrice, bet.Status, bet.OpenedAt, bet.ExpiresAt,
	).Scan(&bet.ID)
}

func (s *Storage) GetBet(ctx context.Context, id int64) (*models.Bet, error) {
	query := `SELECT id, user_id, symbol, direction, stake_amount, opened_price, resolved_price, status, opened_at, expires_at, resolved_at FROM price_bets WHERE id = $1`
	var b models.Bet
	err := s.DB.QueryRowContext(ctx, query, id).Scan(
		&b.ID, &b.UserID, &b.Symbol, &b.Direction, &b.StakeAmount, &b.OpenedPrice, &b.ResolvedPrice, &b.Status, &b.OpenedAt, &b.ExpiresAt, &b.ResolvedAt,
	)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (s *Storage) GetUserBets(ctx context.Context, userID int64) ([]models.Bet, error) {
	query := `SELECT id, user_id, symbol, direction, stake_amount, opened_price, resolved_price, status, opened_at, expires_at, resolved_at FROM price_bets WHERE user_id = $1 ORDER BY opened_at DESC`
	rows, err := s.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bets []models.Bet
	for rows.Next() {
		var b models.Bet
		if err := rows.Scan(&b.ID, &b.UserID, &b.Symbol, &b.Direction, &b.StakeAmount, &b.OpenedPrice, &b.ResolvedPrice, &b.Status, &b.OpenedAt, &b.ExpiresAt, &b.ResolvedAt); err != nil {
			return nil, err
		}
		bets = append(bets, b)
	}
	return bets, nil
}

func (s *Storage) GetExpiredOpenBets(ctx context.Context) ([]models.Bet, error) {
	query := `SELECT id, user_id, symbol, direction, stake_amount, opened_price, resolved_price, status, opened_at, expires_at, resolved_at FROM price_bets WHERE status = 'OPEN' AND expires_at < NOW()`
	rows, err := s.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bets []models.Bet
	for rows.Next() {
		var b models.Bet
		if err := rows.Scan(&b.ID, &b.UserID, &b.Symbol, &b.Direction, &b.StakeAmount, &b.OpenedPrice, &b.ResolvedPrice, &b.Status, &b.OpenedAt, &b.ExpiresAt, &b.ResolvedAt); err != nil {
			return nil, err
		}
		bets = append(bets, b)
	}
	return bets, nil
}

func (s *Storage) UpdateBetStatus(ctx context.Context, bet *models.Bet) error {
	query := `UPDATE price_bets SET status = $1, resolved_price = $2, resolved_at = $3 WHERE id = $4`
	_, err := s.DB.ExecContext(ctx, query, bet.Status, bet.ResolvedPrice, bet.ResolvedAt, bet.ID)
	return err
}

func (s *Storage) ResetBets(ctx context.Context, userID int64) error {
	query := `DELETE FROM price_bets WHERE user_id = $1`
	_, err := s.DB.ExecContext(ctx, query, userID)
	return err
}
