package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"bets-service/internal/models"
	"bets-service/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/shopspring/decimal"
)

type Handler struct {
	Service *service.BetsService
}

func NewHandler(svc *service.BetsService) *Handler {
	return &Handler{Service: svc}
}

func (h *Handler) InitRoutes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Route("/api/v1/bets", func(r chi.Router) {
		r.Post("/", h.createBet)
		r.Get("/", h.getUserBets)
		r.Post("/reset", h.handleResetBets)
	})

	return r
}

type CreateBetRequest struct {
	UserID      int64           `json:"user_id"`
	Symbol      string          `json:"symbol"`
	Direction   string          `json:"direction"`
	StakeAmount decimal.Decimal `json:"stake_amount"`
	DurationSec int             `json:"duration_sec"`
}

func (h *Handler) createBet(w http.ResponseWriter, r *http.Request) {
	var req CreateBetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.UserID == 0 || req.Symbol == "" || req.StakeAmount.LessThanOrEqual(decimal.Zero) || req.DurationSec <= 0 {
		http.Error(w, "Invalid parameters", http.StatusBadRequest)
		return
	}

	var dir models.BetDirection
	if req.Direction == "UP" {
		dir = models.DirectionUp
	} else if req.Direction == "DOWN" {
		dir = models.DirectionDown
	} else {
		http.Error(w, "Invalid direction (UP/DOWN)", http.StatusBadRequest)
		return
	}

	bet, err := h.Service.CreateBet(r.Context(), req.UserID, req.Symbol, dir, req.StakeAmount, req.DurationSec)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bet)
}

func (h *Handler) getUserBets(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		http.Error(w, "Missing user_id", http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}

	bets, err := h.Service.GetUserBets(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bets)
}

func (h *Handler) handleResetBets(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		http.Error(w, "Missing user_id", http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}

	if err := h.Service.ResetBets(r.Context(), userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
