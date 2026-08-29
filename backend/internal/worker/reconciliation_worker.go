package worker

import (
	"context"
	"log/slog"
	"time"
	"github.com/prachii06/LedgerX/internal/repository"
	"github.com/prachii06/LedgerX/internal/services"
)

type ReconciliationWorker struct {
	repository *repository.ReconciliationRepository
	service    *services.ReconciliationService
}

func NewReconciliationWorker(
	repository *repository.ReconciliationRepository,
	service *services.ReconciliationService,
) *ReconciliationWorker {
	return &ReconciliationWorker{
		repository: repository,
		service:    service,
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
			w.reconcileTransactions(ctx)
		}
	}
}


func (w *ReconciliationWorker) reconcileTransactions(
	ctx context.Context,
) {
	transactionIDs, err := w.repository.GetTransactionsNeedingReconciliation(ctx)
	if err != nil {
		slog.Error(
			"failed to get transactions needing reconciliation",
			"error",
			err,
		)
		return
	}

	for _, transactionID := range transactionIDs {
		slog.Info(
			"reconciling transaction",
			"transaction_id",
			transactionID,
		)

		_, err := w.service.ReconcileTransaction(
			ctx,
			transactionID,
		)

		if err != nil {
			slog.Error(
				"failed to reconcile transaction",
				"transaction_id",
				transactionID,
				"error",
				err,
			)
		}
	}
}