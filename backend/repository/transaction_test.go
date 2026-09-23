package repository_test

import (
	"testing"

	"github.com/RR-Sahoo/fintech-starter-kit/backend/models"
	"github.com/RR-Sahoo/fintech-starter-kit/backend/repository"
)

func TestMemoryTransactionRepository(t *testing.T) {
	repo := repository.NewMemoryTransactionRepository()

	t.Run("GetAll returns preseeded data and summary", func(t *testing.T) {
		txs, summary := repo.GetAll("", 0)
		if len(txs) != 10 {
			t.Fatalf("expected 10 initial transactions, got %d", len(txs))
		}
		if summary.TotalCount != 10 {
			t.Errorf("expected summary TotalCount 10, got %d", summary.TotalCount)
		}
		if summary.PendingCount != 5 {
			t.Errorf("expected 5 pending transactions, got %d", summary.PendingCount)
		}
		if summary.SuccessCount != 3 {
			t.Errorf("expected 3 success transactions, got %d", summary.SuccessCount)
		}
		if summary.FailedCount != 2 {
			t.Errorf("expected 2 failed transactions, got %d", summary.FailedCount)
		}
		if summary.SuccessRatePercentage != 30.0 {
			t.Errorf("expected 30%% success rate, got %f", summary.SuccessRatePercentage)
		}
	})

	t.Run("GetByID exists and not exists", func(t *testing.T) {
		tx, err := repo.GetByID("a1029384-b56c-48de-9f12-000000000001")
		if err != nil {
			t.Fatalf("unexpected error finding transaction: %v", err)
		}
		if tx.ID != "a1029384-b56c-48de-9f12-000000000001" {
			t.Errorf("expected ID a1029384-b56c-48de-9f12-000000000001, got %s", tx.ID)
		}

		_, err = repo.GetByID("invalid-id")
		if err != repository.ErrTransactionNotFound {
			t.Errorf("expected ErrTransactionNotFound, got %v", err)
		}
	})

	t.Run("Reconcile pending flips to success", func(t *testing.T) {
		tx, err := repo.Reconcile("a1029384-b56c-48de-9f12-000000000005")
		if err != nil {
			t.Fatalf("unexpected error reconciling: %v", err)
		}
		if tx.Status != models.StatusSuccess {
			t.Errorf("expected status SUCCESS, got %s", tx.Status)
		}

		// Trying to reconcile again should return ErrTransactionNotPending
		_, err = repo.Reconcile("a1029384-b56c-48de-9f12-000000000005")
		if err != repository.ErrTransactionNotPending {
			t.Errorf("expected ErrTransactionNotPending, got %v", err)
		}
	})
}
