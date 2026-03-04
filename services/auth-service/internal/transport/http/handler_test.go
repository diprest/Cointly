package http

import (
	"auth-service/internal/models"
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Register(ctx context.Context, login, password string) (*models.User, error) {
	args := m.Called(ctx, login, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockAuthService) Login(ctx context.Context, login, password string) (string, int, error) {
	args := m.Called(ctx, login, password)
	return args.String(0), args.Int(1), args.Error(2)
}

func (m *MockAuthService) ValidateToken(tokenString string) (int, error) {
	args := m.Called(tokenString)
	return args.Int(0), args.Error(1)
}

func (m *MockAuthService) ChangeLogin(ctx context.Context, userID int, newLogin string) error {
	args := m.Called(ctx, userID, newLogin)
	return args.Error(0)
}

func (m *MockAuthService) ChangePassword(ctx context.Context, userID int, oldPassword, newPassword string) error {
	args := m.Called(ctx, userID, oldPassword, newPassword)
	return args.Error(0)
}

func TestHandler_Register(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := new(MockAuthService)
		handler := NewHandler(mockService)
		router := handler.InitRoutes()

		user := &models.User{ID: 1, Login: "test", CreatedAt: time.Now()}
		mockService.On("Register", mock.Anything, "test", "pass").Return(user, nil)

		body := bytes.NewBufferString(`{"login":"test", "password":"pass"}`)
		req, _ := http.NewRequest("POST", "/api/v1/auth/register", body)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("BadRequest", func(t *testing.T) {
		mockService := new(MockAuthService)
		handler := NewHandler(mockService)
		router := handler.InitRoutes()

		body := bytes.NewBufferString(`{"login":"test"}`)
		req, _ := http.NewRequest("POST", "/api/v1/auth/register", body)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("InternalError", func(t *testing.T) {
		mockService := new(MockAuthService)
		handler := NewHandler(mockService)
		router := handler.InitRoutes()

		mockService.On("Register", mock.Anything, "test", "pass").Return(nil, errors.New("db error"))

		body := bytes.NewBufferString(`{"login":"test", "password":"pass"}`)
		req, _ := http.NewRequest("POST", "/api/v1/auth/register", body)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestHandler_Login(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := new(MockAuthService)
		handler := NewHandler(mockService)
		router := handler.InitRoutes()

		mockService.On("Login", mock.Anything, "test", "pass").Return("token123", 1, nil)

		body := bytes.NewBufferString(`{"login":"test", "password":"pass"}`)
		req, _ := http.NewRequest("POST", "/api/v1/auth/login", body)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "token123")
	})

	t.Run("Unauthorized", func(t *testing.T) {
		mockService := new(MockAuthService)
		handler := NewHandler(mockService)
		router := handler.InitRoutes()

		mockService.On("Login", mock.Anything, "test", "pass").Return("", 0, errors.New("invalid credentials"))

		body := bytes.NewBufferString(`{"login":"test", "password":"pass"}`)
		req, _ := http.NewRequest("POST", "/api/v1/auth/login", body)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestHandler_Validate(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := new(MockAuthService)
		handler := NewHandler(mockService)
		router := handler.InitRoutes()

		mockService.On("ValidateToken", "valid_token").Return(1, nil)

		req, _ := http.NewRequest("GET", "/api/v1/auth/validate", nil)
		req.Header.Set("Authorization", "Bearer valid_token")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"user_id":1`)
	})

	t.Run("NoHeader", func(t *testing.T) {
		mockService := new(MockAuthService)
		handler := NewHandler(mockService)
		router := handler.InitRoutes()

		req, _ := http.NewRequest("GET", "/api/v1/auth/validate", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("InvalidHeaderFormat", func(t *testing.T) {
		mockService := new(MockAuthService)
		handler := NewHandler(mockService)
		router := handler.InitRoutes()

		req, _ := http.NewRequest("GET", "/api/v1/auth/validate", nil)
		req.Header.Set("Authorization", "Token valid_token")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("InvalidToken", func(t *testing.T) {
		mockService := new(MockAuthService)
		handler := NewHandler(mockService)
		router := handler.InitRoutes()

		mockService.On("ValidateToken", "invalid_token").Return(0, errors.New("invalid token"))

		req, _ := http.NewRequest("GET", "/api/v1/auth/validate", nil)
		req.Header.Set("Authorization", "Bearer invalid_token")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
