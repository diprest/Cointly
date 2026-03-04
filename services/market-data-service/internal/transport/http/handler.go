package http

import (
	"encoding/json"
	"log/slog"
	"market-data-service/internal/service"
	"net/http"
)

type Handler struct {
	service *service.MarketService
}

func NewHandler(service *service.MarketService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/market/symbols", h.handleGetSymbols)
	mux.HandleFunc("GET /api/v1/market/ticker", h.handleGetTicker)
}

func (h *Handler) handleGetSymbols(w http.ResponseWriter, r *http.Request) {
	symbols, err := h.service.GetAllSymbols(r.Context())
	if err != nil {
		slog.Error("Failed to get symbols", "error", err)
		http.Error(w, `{"error": "internal server error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(symbols)
}

func (h *Handler) handleGetTicker(w http.ResponseWriter, r *http.Request) {
	symbol := r.URL.Query().Get("symbol")
	if symbol == "" {
		http.Error(w, `{"error": "symbol param is required"}`, http.StatusBadRequest)
		return
	}

	ticker, err := h.service.GetTicker(r.Context(), symbol)
	if err != nil {
		slog.Warn("Ticker not found", "symbol", symbol, "error", err)
		http.Error(w, `{"error": "price not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ticker)
}
