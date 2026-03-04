package http

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"trading-service/internal/models"
	"trading-service/internal/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *service.TradingService
}

func NewHandler(svc *service.TradingService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) InitRoutes() *gin.Engine {
	r := gin.Default()
	api := r.Group("/api/v1")
	{
		api.POST("/trading/orders", h.createOrder)
		api.DELETE("/trading/orders/:id", h.cancelOrder)
		api.GET("/trading/orders", h.getOrders)
		api.POST("/trading/reset", h.handleResetOrders)
	}
	return r
}

func (h *Handler) createOrder(c *gin.Context) {
	var o models.Order
	if err := c.ShouldBindJSON(&o); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if o.UserID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	if err := h.svc.CreateOrder(&o); err != nil {
		log.Printf("Failed to create order: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, o)
}

func (h *Handler) cancelOrder(c *gin.Context) {
	idStr := strings.TrimSpace(c.Param("id"))
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}

	userIDStr := c.Query("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil || userID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	if err := h.svc.CancelOrder(c.Request.Context(), uint(id), int64(userID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "cancelled", "order_id": id})
}

func (h *Handler) getOrders(c *gin.Context) {
	userIDStr := c.Query("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil || userID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	orders, err := h.svc.GetUserOrders(int64(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}

func (h *Handler) handleResetOrders(c *gin.Context) {
	userIDStr := c.Query("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil || userID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	if err := h.svc.ResetOrders(c.Request.Context(), int64(userID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
