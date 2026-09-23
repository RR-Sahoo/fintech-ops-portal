package repository

import (
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/RR-Sahoo/fintech-starter-kit/backend/models"
)

var (
	ErrTransactionNotFound     = errors.New("transaction not found")
	ErrTransactionNotPending   = errors.New("transaction is not in PENDING status")
	ErrInvalidTransactionStatus = errors.New("invalid transaction status")
)

// TransactionSummary represents aggregated metrics over transactions in USD cents.
type TransactionSummary struct {
	TotalVolumeCents      int64   `json:"total_volume_cents"`
	TotalVolumeFormatted  string  `json:"total_volume_formatted"`
	TotalCount            int     `json:"total_count"`
	PendingCount          int     `json:"pending_count"`
	SuccessCount          int     `json:"success_count"`
	FailedCount           int     `json:"failed_count"`
	SuccessRatePercentage float64 `json:"success_rate_percentage"`
}

// TransactionRepository defines the data access contract.
type TransactionRepository interface {
	GetAll(statusFilter string, limit int) ([]models.Transaction, TransactionSummary)
	GetByID(id string) (*models.Transaction, error)
	Reconcile(id string) (*models.Transaction, error)
}

// MemoryTransactionRepository is an in-memory thread-safe implementation of TransactionRepository.
type MemoryTransactionRepository struct {
	mu           sync.RWMutex
	transactions []models.Transaction
}

// NewMemoryTransactionRepository initializes the in-memory repository pre-seeded with realistic mock fintech data in USD.
func NewMemoryTransactionRepository() *MemoryTransactionRepository {
	now := time.Now().UTC()

	initialTransactions := []models.Transaction{
		{
			ID:              "a1029384-b56c-48de-9f12-000000000001",
			UserID:          "usr_fin_9821",
			AmountCents:     12500, // $125.00
			Currency:        "USD",
			Status:          models.StatusPending,
			ReferenceID:     "ACH/428910293812/PAY_MERCHANT",
			Type:            models.TypeDebit,
			CreatedAt:       now.Add(-2 * time.Hour),
			FormattedAmount: models.FormatUSD(12500),
		},
		{
			ID:              "a1029384-b56c-48de-9f12-000000000002",
			UserID:          "usr_fin_4412",
			AmountCents:     450000, // $4,500.00
			Currency:        "USD",
			Status:          models.StatusSuccess,
			ReferenceID:     "WIRE/W90281203810/SALARY_CREDIT",
			Type:            models.TypeCredit,
			CreatedAt:       now.Add(-4 * time.Hour),
			FormattedAmount: models.FormatUSD(450000),
		},
		{
			ID:              "a1029384-b56c-48de-9f12-000000000003",
			UserID:          "usr_fin_1092",
			AmountCents:     125000, // $1,250.00
			Currency:        "USD",
			Status:          models.StatusPending,
			ReferenceID:     "STRIPE/TXN_881920384910",
			Type:            models.TypeDebit,
			CreatedAt:       now.Add(-30 * time.Minute),
			FormattedAmount: models.FormatUSD(125000),
		},
		{
			ID:              "a1029384-b56c-48de-9f12-000000000004",
			UserID:          "usr_fin_3389",
			AmountCents:     7500, // $75.00
			Currency:        "USD",
			Status:          models.StatusFailed,
			ReferenceID:     "CARD/428919920192/BILL_UTILITY",
			Type:            models.TypeDebit,
			CreatedAt:       now.Add(-1 * time.Hour),
			FormattedAmount: models.FormatUSD(7500),
		},
		{
			ID:              "a1029384-b56c-48de-9f12-000000000005",
			UserID:          "usr_fin_9821",
			AmountCents:     50000, // $500.00
			Currency:        "USD",
			Status:          models.StatusPending,
			ReferenceID:     "ACH/2026092400192/VENDOR_PAY",
			Type:            models.TypeDebit,
			CreatedAt:       now.Add(-15 * time.Minute),
			FormattedAmount: models.FormatUSD(50000),
		},
		{
			ID:              "a1029384-b56c-48de-9f12-000000000006",
			UserID:          "usr_fin_7718",
			AmountCents:     22000, // $220.00
			Currency:        "USD",
			Status:          models.StatusPending,
			ReferenceID:     "REFUND/428901928374/REFUND_ORDER",
			Type:            models.TypeCredit,
			CreatedAt:       now.Add(-5 * time.Hour),
			FormattedAmount: models.FormatUSD(22000),
		},
		{
			ID:              "a1029384-b56c-48de-9f12-000000000007",
			UserID:          "usr_fin_5521",
			AmountCents:     320000, // $3,200.00
			Currency:        "USD",
			Status:          models.StatusPending,
			ReferenceID:     "FEDWIRE/F202609240092/EQUITY_DEP",
			Type:            models.TypeCredit,
			CreatedAt:       now.Add(-10 * time.Minute),
			FormattedAmount: models.FormatUSD(320000),
		},
		{
			ID:              "a1029384-b56c-48de-9f12-000000000008",
			UserID:          "usr_fin_4412",
			AmountCents:     4999, // $49.99
			Currency:        "USD",
			Status:          models.StatusFailed,
			ReferenceID:     "CARD/AUTH_8910283019",
			Type:            models.TypeDebit,
			CreatedAt:       now.Add(-6 * time.Hour),
			FormattedAmount: models.FormatUSD(4999),
		},
		{
			ID:              "a1029384-b56c-48de-9f12-000000000009",
			UserID:          "usr_fin_6634",
			AmountCents:     85000, // $850.00
			Currency:        "USD",
			Status:          models.StatusSuccess,
			ReferenceID:     "ACH/428900192834/LOAN_PAYMENT",
			Type:            models.TypeDebit,
			CreatedAt:       now.Add(-8 * time.Hour),
			FormattedAmount: models.FormatUSD(85000),
		},
		{
			ID:              "a1029384-b56c-48de-9f12-000000000010",
			UserID:          "usr_fin_2209",
			AmountCents:     18750, // $187.50
			Currency:        "USD",
			Status:          models.StatusSuccess,
			ReferenceID:     "CARD/428938491029/SAAS_SUB",
			Type:            models.TypeDebit,
			CreatedAt:       now.Add(-12 * time.Hour),
			FormattedAmount: models.FormatUSD(18750),
		},
	}

	return &MemoryTransactionRepository{
		transactions: initialTransactions,
	}
}

// GetAll returns transactions optionally filtered by status with limit and calculated summary metrics.
func (r *MemoryTransactionRepository) GetAll(statusFilter string, limit int) ([]models.Transaction, TransactionSummary) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var filtered []models.Transaction
	var totalVolumeCents int64
	var pendingCount, successCount, failedCount int

	targetStatus := strings.ToUpper(strings.TrimSpace(statusFilter))

	for _, tx := range r.transactions {
		// Calculate summary statistics across all transactions in store
		totalVolumeCents += tx.AmountCents
		switch tx.Status {
		case models.StatusPending:
			pendingCount++
		case models.StatusSuccess:
			successCount++
		case models.StatusFailed:
			failedCount++
		}

		// Apply filter if specified
		if targetStatus == "" || tx.Status == targetStatus {
			filtered = append(filtered, tx)
		}
	}

	totalCount := len(r.transactions)
	var successRate float64
	if totalCount > 0 {
		successRate = (float64(successCount) / float64(totalCount)) * 100.0
	}

	summary := TransactionSummary{
		TotalVolumeCents:      totalVolumeCents,
		TotalVolumeFormatted:  models.FormatUSD(totalVolumeCents),
		TotalCount:            totalCount,
		PendingCount:          pendingCount,
		SuccessCount:          successCount,
		FailedCount:           failedCount,
		SuccessRatePercentage: successRate,
	}

	if limit > 0 && len(filtered) > limit {
		filtered = filtered[:limit]
	}

	return filtered, summary
}

// GetByID searches for a transaction by ID.
func (r *MemoryTransactionRepository) GetByID(id string) (*models.Transaction, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, tx := range r.transactions {
		if tx.ID == id {
			res := tx
			return &res, nil
		}
	}
	return nil, ErrTransactionNotFound
}

// Reconcile flips a PENDING transaction to SUCCESS and updates its timestamp.
func (r *MemoryTransactionRepository) Reconcile(id string) (*models.Transaction, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i := range r.transactions {
		if r.transactions[i].ID == id {
			if r.transactions[i].Status != models.StatusPending {
				return nil, ErrTransactionNotPending
			}

			now := time.Now().UTC()
			r.transactions[i].Status = models.StatusSuccess
			r.transactions[i].UpdatedAt = now
			res := r.transactions[i]
			return &res, nil
		}
	}

	return nil, ErrTransactionNotFound
}
