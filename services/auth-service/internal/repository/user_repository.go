package repository

import (
	"auth-service/internal/models"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type DBPool interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
}

type UserRepository struct {
	db DBPool
}

func NewUserRepository(db DBPool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	query := `INSERT INTO users (login, password_hash) VALUES ($1, $2) RETURNING id, created_at`
	err := r.db.QueryRow(ctx, query, user.Login, user.PasswordHash).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return fmt.Errorf("user already exists")
		}
		return err
	}
	return nil
}

func (r *UserRepository) GetUserByLogin(ctx context.Context, login string) (*models.User, error) {
	query := `SELECT id, login, password_hash, created_at FROM users WHERE login = $1`
	var user models.User
	err := r.db.QueryRow(ctx, query, login).Scan(&user.ID, &user.Login, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	query := `SELECT id, login, password_hash, created_at FROM users WHERE id = $1`
	var user models.User
	err := r.db.QueryRow(ctx, query, id).Scan(&user.ID, &user.Login, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) UpdateLogin(ctx context.Context, userID int, newLogin string) error {
	query := `UPDATE users SET login = $1 WHERE id = $2`
	commandTag, err := r.db.Exec(ctx, query, newLogin, userID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return fmt.Errorf("login already taken")
		}
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return errors.New("user not found")
	}
	return nil
}

func (r *UserRepository) UpdatePassword(ctx context.Context, userID int, newPasswordHash string) error {
	query := `UPDATE users SET password_hash = $1 WHERE id = $2`
	commandTag, err := r.db.Exec(ctx, query, newPasswordHash, userID)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return errors.New("user not found")
	}
	return nil
}
