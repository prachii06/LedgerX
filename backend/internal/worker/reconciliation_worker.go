package worker

import (
	"context"
	"log/slog"
	"time"
	"github.com/prachii06/LedgerX/internal/services"
)

type ReconciliationWorker struct {
	service *services.ReconciliationService
}

func NewReconciliationWorker(
	service *services.ReconciliationService,
) *ReconciliationWorker {
	return &ReconciliationWorker{
		service: service,
	}
}

func (w *ReconciliationWorker) Start(
	ctx context.Context,
	interval time.Duration,
) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			slog.Info("reconciliation worker stopped")
			return

		case <-ticker.C:
			slog.Info("reconciliation worker running")
		}
	}
}