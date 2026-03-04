package worker

import (
	"context"
	"log/slog"
	"time"

	"bets-service/internal/service"
)

type Resolver struct {
	Service *service.BetsService
}

func NewResolver(svc *service.BetsService) *Resolver {
	return &Resolver{Service: svc}
}

func (r *Resolver) Start(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	slog.Info("Starting Bet Resolver Worker")

	for {
		select {
		case <-ctx.Done():
			slog.Info("Stopping Bet Resolver Worker")
			return
		case <-ticker.C:
			r.processExpiredBets(ctx)
		}
	}
}

func (r *Resolver) processExpiredBets(ctx context.Context) {
	bets, err := r.Service.Repo.GetExpiredOpenBets(ctx)
	if err != nil {
		slog.Error("Failed to get expired bets", "error", err)
		return
	}

	for _, bet := range bets {
		resolveCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		if err := r.Service.ResolveBet(resolveCtx, &bet); err != nil {
			slog.Error("Failed to resolve bet", "bet_id", bet.ID, "error", err)
		} else {
			slog.Info("Bet resolved", "bet_id", bet.ID, "status", bet.Status, "win", bet.Status == "WON")
		}
	}
}
