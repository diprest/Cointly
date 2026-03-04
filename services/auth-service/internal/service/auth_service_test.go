package service

import (
	"auth-service/internal/models"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetUserByLogin(ctx context.Context, login string) (*models.User, error) {
	args := m.Called(ctx, login)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) UpdateLogin(ctx context.Context, userID int, newLogin string) error {
	args := m.Called(ctx, userID, newLogin)
	return args.Error(0)
}

func (m *MockUserRepository) UpdatePassword(ctx context.Context, userID int, newPasswordHash string) error {
	args := m.Called(ctx, userID, newPasswordHash)
	return args.Error(0)
}

func TestAuthService_Register(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewAuthService(mockRepo, "secret")

	ctx := context.Background()
	login := "testuser"
	password := "password"

	mockRepo.On("CreateUser", ctx, mock.AnythingOfType("*models.User")).Return(nil)

	user, err := service.Register(ctx, login, password)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, login, user.Login)
	assert.NotEmpty(t, user.PasswordHash)

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestAuthService_Login_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewAuthService(mockRepo, "secret")

	ctx := context.Background()
	login := "testuser"
	password := "password"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	existingUser := &models.User{
		ID:           1,
		Login:        login,
		PasswordHash: string(hashedPassword),
	}

	mockRepo.On("GetUserByLogin", ctx, login).Return(existingUser, nil)

	token, userID, err := service.Login(ctx, login, password)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.Equal(t, 1, userID)
	mockRepo.AssertExpectations(t)
}

func TestAuthService_Login_InvalidPassword(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewAuthService(mockRepo, "secret")

	ctx := context.Background()
	login := "testuser"
	password := "password"
	wrongPassword := "wrongpassword"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	existingUser := &models.User{
		ID:           1,
		Login:        login,
		PasswordHash: string(hashedPassword),
	}

	mockRepo.On("GetUserByLogin", ctx, login).Return(existingUser, nil)

	token, _, err := service.Login(ctx, login, wrongPassword)

	assert.Error(t, err)
	assert.Equal(t, "invalid credentials", err.Error())
	assert.Empty(t, token)
	mockRepo.AssertExpectations(t)
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewAuthService(mockRepo, "secret")

	ctx := context.Background()
	login := "nonexistent"
	password := "password"

	mockRepo.On("GetUserByLogin", ctx, login).Return(nil, nil)

	token, _, err := service.Login(ctx, login, password)

	assert.Error(t, err)
	assert.Equal(t, "invalid credentials", err.Error())
	assert.Empty(t, token)
	mockRepo.AssertExpectations(t)
}

func TestAuthService_ValidateToken_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewAuthService(mockRepo, "secret")

	ctx := context.Background()
	login := "testuser"
	password := "password"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	existingUser := &models.User{ID: 1, Login: login, PasswordHash: string(hashedPassword)}

	mockRepo.On("GetUserByLogin", ctx, login).Return(existingUser, nil)
	tokenString, _, _ := service.Login(ctx, login, password)

	userID, err := service.ValidateToken(tokenString)

	assert.NoError(t, err)
	assert.Equal(t, 1, userID)
}

func TestAuthService_ValidateToken_Invalid(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewAuthService(mockRepo, "secret")

	userID, err := service.ValidateToken("invalid.token.string")

	assert.Error(t, err)
	assert.Equal(t, 0, userID)
}
