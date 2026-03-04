package repository

import (
	"context"
	"testing"
	"time"

	"auth-service/internal/models"

	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository_CreateUser(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()

	repo := NewUserRepository(mock)

	user := &models.User{
		Login:        "testuser",
		PasswordHash: "hashedpassword",
	}

	rows := pgxmock.NewRows([]string{"id", "created_at"}).AddRow(1, time.Now())
	mock.ExpectQuery("INSERT INTO users").
		WithArgs(user.Login, user.PasswordHash).
		WillReturnRows(rows)

	err = repo.CreateUser(context.Background(), user)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserRepository_GetUserByLogin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()

	repo := NewUserRepository(mock)

	expectedUser := &models.User{
		ID:           1,
		Login:        "testuser",
		PasswordHash: "hashedpassword",
		CreatedAt:    time.Now(),
	}

	rows := pgxmock.NewRows([]string{"id", "login", "password_hash", "created_at"}).
		AddRow(expectedUser.ID, expectedUser.Login, expectedUser.PasswordHash, expectedUser.CreatedAt)

	mock.ExpectQuery("SELECT id, login, password_hash, created_at FROM users WHERE login = \\$1").
		WithArgs(expectedUser.Login).
		WillReturnRows(rows)

	user, err := repo.GetUserByLogin(context.Background(), expectedUser.Login)
	assert.NoError(t, err)
	assert.Equal(t, expectedUser.ID, user.ID)
	assert.Equal(t, expectedUser.Login, user.Login)
	assert.Equal(t, expectedUser.PasswordHash, user.PasswordHash)
	assert.WithinDuration(t, expectedUser.CreatedAt, user.CreatedAt, time.Second)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserRepository_GetUserByID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()

	repo := NewUserRepository(mock)

	expectedUser := &models.User{
		ID:           1,
		Login:        "testuser",
		PasswordHash: "hashedpassword",
		CreatedAt:    time.Now(),
	}

	rows := pgxmock.NewRows([]string{"id", "login", "password_hash", "created_at"}).
		AddRow(expectedUser.ID, expectedUser.Login, expectedUser.PasswordHash, expectedUser.CreatedAt)

	mock.ExpectQuery("SELECT id, login, password_hash, created_at FROM users WHERE id = \\$1").
		WithArgs(expectedUser.ID).
		WillReturnRows(rows)

	user, err := repo.GetUserByID(context.Background(), expectedUser.ID)
	assert.NoError(t, err)
	assert.Equal(t, expectedUser.ID, user.ID)
	assert.Equal(t, expectedUser.Login, user.Login)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
