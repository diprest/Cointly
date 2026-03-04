package http

import (
	"auth-service/internal/models"
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthService interface {
	Register(ctx context.Context, login, password string) (*models.User, error)
	Login(ctx context.Context, login, password string) (string, int, error)
	ValidateToken(tokenString string) (int, error)
	ChangePassword(ctx context.Context, userID int, oldPassword, newPassword string) error
	ChangeLogin(ctx context.Context, userID int, newLogin string) error
}

type Handler struct {
	authService AuthService
}

func NewHandler(authService AuthService) *Handler {
	return &Handler{authService: authService}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.Default()

	api := router.Group("/api/v1/auth")
	{
		api.POST("/register", h.register)
		api.POST("/login", h.login)
		api.GET("/validate", h.validate)

		protected := api.Group("/user")
		protected.Use(h.userIdentity)
		protected.POST("/change-login", h.changeLogin)
		protected.POST("/change-password", h.changePassword)
	}

	return router
}

func (h *Handler) userIdentity(c *gin.Context) {
	header := c.GetHeader("Authorization")
	if header == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "empty auth header"})
		return
	}

	parts := strings.Split(header, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid auth header"})
		return
	}

	userID, err := h.authService.ValidateToken(parts[1])
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.Set("userID", userID)
}

func (h *Handler) register(c *gin.Context) {
	var input models.RegisterRequest
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.authService.Register(c.Request.Context(), input.Login, input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *Handler) login(c *gin.Context) {
	var input models.LoginRequest
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, userID, err := h.authService.Login(c.Request.Context(), input.Login, input.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.AuthResponse{Token: token, UserID: userID})
}

func (h *Handler) validate(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "empty auth header"})
		return
	}

	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid auth header"})
		return
	}

	userID, err := h.authService.ValidateToken(headerParts[1])
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user_id": userID})
}

func (h *Handler) changeLogin(c *gin.Context) {
	userID, _ := c.Get("userID")

	var input models.ChangeLoginRequest
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.authService.ChangeLogin(c.Request.Context(), userID.(int), input.NewLogin); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *Handler) changePassword(c *gin.Context) {
	userID, _ := c.Get("userID")

	var input models.ChangePasswordRequest
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.authService.ChangePassword(c.Request.Context(), userID.(int), input.OldPassword, input.NewPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
