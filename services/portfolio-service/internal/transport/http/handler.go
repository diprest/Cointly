package http

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"portfolio-service/internal/models"
	"portfolio-service/internal/service"
	"strconv"

	"github.com/shopspring/decimal"
)

type Handler struct {
	service *service.PortfolioService
}

func NewHandler(svc *service.PortfolioService) *Handler {
	return &Handler{service: svc}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/portfolio/balance", h.handleGetBalance)
	mux.HandleFunc("/api/v1/portfolio/list", h.handleGetPortfolio)
	mux.HandleFunc("/api/v1/portfolio/lock", h.handleLockFunds)
	mux.HandleFunc("/api/v1/portfolio/unlock", h.handleUnlockFunds)
	mux.HandleFunc("/api/v1/portfolio/transfer", h.handleTransferFunds)
	mux.HandleFunc("/api/v1/portfolio/balance/change", h.handleChangeBalance)
	mux.HandleFunc("/api/v1/portfolio/reset", h.handleResetBalance)
}

func (h *Handler) handleGetBalance(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	asset := r.URL.Query().Get("asset")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil || userID <= 0 || asset == "" {
		http.Error(w, `{"error": "user_id and asset required"}`, http.StatusBadRequest)
		return
	}

	bal, err := h.service.GetBalance(r.Context(), userID, asset)
	if err != nil {
		slog.Error("GetBalance error", "err", err)
		http.Error(w, `{"error": "internal error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id":    bal.UserID,
		"asset":      bal.Asset,
		"amount":     bal.Amount.String(),
		"locked_bal": bal.LockedBal.String(),
	})
}

func (h *Handler) handleGetPortfolio(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil || userID <= 0 {
		http.Error(w, `{"error": "user_id required"}`, http.StatusBadRequest)
		return
	}

	balances, err := h.service.GetPortfolio(r.Context(), userID)
	if err != nil {
		slog.Error("GetPortfolio error", "err", err)
		http.Error(w, `{"error": "internal error"}`, http.StatusInternalServerError)
		return
	}

	resp := make([]map[string]interface{}, len(balances))
	for i, b := range balances {
		resp[i] = map[string]interface{}{
			"user_id":    b.UserID,
			"asset":      b.Asset,
			"amount":     b.Amount.String(),
			"locked_bal": b.LockedBal.String(),
			"total_cost": b.TotalCost.String(),
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) handleLockFunds(w http.ResponseWriter, r *http.Request) {
	var req models.LockUnlockRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	amount, _ := decimal.NewFromString(req.Amount)

	if err := h.service.LockFunds(r.Context(), req.UserID, req.Asset, amount); err != nil {
		if errors.Is(err, service.ErrInsufficientFunds) {
			http.Error(w, `{"error": "insufficient funds"}`, http.StatusForbidden)
			return
		}
		slog.Error("Lock error", "err", err)
		http.Error(w, `{"error": "internal error"}`, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) handleUnlockFunds(w http.ResponseWriter, r *http.Request) {
	var req models.LockUnlockRequest
	json.NewDecoder(r.Body).Decode(&req)
	amount, _ := decimal.NewFromString(req.Amount)

	if err := h.service.UnlockFunds(r.Context(), req.UserID, req.Asset, amount); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) handleTransferFunds(w http.ResponseWriter, r *http.Request) {
	var req models.TransferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	amt, _ := decimal.NewFromString(req.Amount)
	cost, _ := decimal.NewFromString(req.Cost)

	err := h.service.TransferFunds(r.Context(), req.UserID, req.Asset, amt, cost, req.Side)
	if err != nil {
		slog.Error("Transfer error", "err", err)
		http.Error(w, `{"error": "transfer failed"}`, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) handleChangeBalance(w http.ResponseWriter, r *http.Request) {
	var req models.LockUnlockRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	amount, _ := decimal.NewFromString(req.Amount)

	if err := h.service.ChangeBalance(r.Context(), req.UserID, req.Asset, amount); err != nil {
		slog.Error("ChangeBalance error", "err", err)
		http.Error(w, `{"error": "change balance failed"}`, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) handleResetBalance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userIDStr := r.URL.Query().Get("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil || userID <= 0 {
		http.Error(w, `{"error": "user_id required"}`, http.StatusBadRequest)
		return
	}

	if err := h.service.ResetBalance(r.Context(), userID); err != nil {
		slog.Error("ResetBalance error", "err", err)
		http.Error(w, `{"error": "reset failed"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
