package storage

import (
	"database/sql"
	"errors"
	"time"
	"trading-service/internal/models"

	_ "github.com/lib/pq"
)

type OrderRepository interface {
	CreateOrder(o *models.Order) error
	GetOrder(id uint) (*models.Order, error)
	UpdateOrderStatus(id uint, status string) error
	GetUserOrders(userID int64) ([]models.Order, error)
	GetActiveOrders() ([]models.Order, error)
	ResetOrders(userID int64) error
}

type Storage struct {
	DB *sql.DB
}

func NewStorage(dsn string) (*Storage, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	query := `
	CREATE TABLE IF NOT EXISTS orders (
		id SERIAL PRIMARY KEY,
		user_id BIGINT NOT NULL,
		symbol VARCHAR(20) NOT NULL,
		side VARCHAR(10) NOT NULL,
		type VARCHAR(10) NOT NULL,
		price DECIMAL(20, 8) NOT NULL,
		amount DECIMAL(20, 8) NOT NULL,
		quote_amount DECIMAL(20, 8) DEFAULT 0,
		status VARCHAR(20) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`
	if _, err := db.Exec(query); err != nil {
		return nil, err
	}

	return &Storage{DB: db}, nil
}

func (s *Storage) CreateOrder(o *models.Order) error {
	query := `
		INSERT INTO orders (user_id, symbol, side, type, price, amount, quote_amount, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id
	`
	err := s.DB.QueryRow(
		query,
		o.UserID, o.Symbol, o.Side, o.Type, o.Price, o.Amount, o.QuoteAmount, o.Status, time.Now(), time.Now(),
	).Scan(&o.ID)

	return err
}

func (s *Storage) GetOrder(id uint) (*models.Order, error) {
	var o models.Order
	query := `SELECT id, user_id, symbol, side, type, price, amount, quote_amount, status FROM orders WHERE id = $1`

	err := s.DB.QueryRow(query, id).Scan(
		&o.ID, &o.UserID, &o.Symbol, &o.Side, &o.Type, &o.Price, &o.Amount, &o.QuoteAmount, &o.Status,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New("order not found")
	}
	return &o, err
}

func (s *Storage) UpdateOrderStatus(id uint, status string) error {
	query := `UPDATE orders SET status = $1, updated_at = $2 WHERE id = $3`
	_, err := s.DB.Exec(query, status, time.Now(), id)
	return err
}

func (s *Storage) GetUserOrders(userID int64) ([]models.Order, error) {
	rows, err := s.DB.Query(`SELECT id, user_id, symbol, side, type, price, amount, quote_amount, status FROM orders WHERE user_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var o models.Order
		if err := rows.Scan(&o.ID, &o.UserID, &o.Symbol, &o.Side, &o.Type, &o.Price, &o.Amount, &o.QuoteAmount, &o.Status); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}

func (s *Storage) GetActiveOrders() ([]models.Order, error) {
	rows, err := s.DB.Query(`SELECT id, user_id, symbol, side, type, price, amount, quote_amount, status FROM orders WHERE status = 'NEW'`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var o models.Order
		if err := rows.Scan(&o.ID, &o.UserID, &o.Symbol, &o.Side, &o.Type, &o.Price, &o.Amount, &o.QuoteAmount, &o.Status); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}

func (s *Storage) ResetOrders(userID int64) error {
	query := `DELETE FROM orders WHERE user_id = $1`
	_, err := s.DB.Exec(query, userID)
	return err
}
