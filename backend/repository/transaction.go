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

// TransactionSummary represents aggregated metrics over transactions.
type TransactionSummary struct {
	TotalVolumePaise       int64   `json:"total_volume_paise"`
	TotalVolumeFormatted   string  `json:"total_volume_formatted"`
	TotalCount             int     `json:"total_count"`
	PendingCount           int     `json:"pending_count"`
	SuccessCount           int     `json:"success_count"`
	FailedCount            int     `json:"failed_count"`
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

// NewMemoryTransactionRepository initializes the in-memory repository pre-seeded with realistic mock fintech data.
func NewMemoryTransactionRepository() *MemoryTransactionRepository {
	now := time.Now().UTC()

	initialTransactions := []models.Transaction{
		{
			ID:              "a1029384-b56c-48de-9f12-000000000001",
			UserID:          "usr_fin_9821",
			AmountPaise:     125050, // ₹1,250.50
			Currency:        "INR",
			Status:          models.StatusPending,
			ReferenceID:     "UPI/428910293812/PAY_MERCHANT",
			Type:            models.TypeDebit,
			CreatedAt:       now.Add(-2 * time.Hour),
			FormattedAmount: models.FormatRupees(125050),
		},
		{
			ID:              "a1029384-b56c-48de-9f12-000000000002",
			UserID:          "usr_fin_4412",
			AmountPaise:     5000000, // ₹50,000.00
			Currency:        "INR",
			Status:          models.StatusSuccess,
			ReferenceID:     "NEFT/N90281203810/SALARY_CREDIT",
			Type:            models.TypeCredit,
			CreatedAt:       now.Add(-4 * time.Hour),
			FormattedAmount: models.FormatRupees(5000000),
		},
		{
			ID:              "a1029384-b56c-48de-9f12-000000000003",
			UserID:          "usr_fin_1092",
			AmountPaise:     349900, // ₹3,499.00
			Currency:        "INR",
			Status:          models.StatusPending,
			ReferenceID:     "PG/TXN_881920384910",
			Type:            models.TypeDebit,
			CreatedAt:       now.Add(-30 * time.Minute),
			FormattedAmount: models.FormatRupees(349900),
		},
		{
			ID:              "a1029384-b56c-48de-9f12-000000000004",
			UserID:          "usr_fin_3389",
			AmountPaise:     75000, // ₹750.00
			Currency:        "INR",
			Status:          models.StatusFailed,
			ReferenceID:     "UPI/428919920192/BILL_UTILITY",
			Type:            models.TypeDebit,
			CreatedAt:       now.Add(-1 * time.Hour),
			FormattedAmount: models.FormatRupees(75000),
		},
		{
			ID:              "a1029384-b56c-48de-9f12-000000000005",
			UserID:          "usr_fin_9821",
			AmountPaise:     1500000, // ₹15,000.00
			Currency:        "INR",
			Status:          models.StatusPending,
			ReferenceID:     "IMPS/2026092400192/VENDOR_PAY",
			Type:            models.TypeDebit,
			CreatedAt:       now.Add(-15 * time.Minute),
			FormattedAmount: models.FormatRupees(1500000),
		},
		{
			ID:              "a1029384-b56c-48de-9f12-000000000006",
			UserID:          "usr_fin_7718",
			AmountPaise:     220000, // ₹2,200.00
			Currency:        "INR",
			Status:          models.StatusPending,
			ReferenceID:     "UPI/428901928374/REFUND_ORDER",
			Type:            models.TypeCredit,
			CreatedAt:       now.Add(-5 * time.Hour),
			FormattedAmount: models.FormatRupees(220000),
		},
		{
			ID:              "a1029384-b56c-48de-9f12-000000000007",
			UserID:          "usr_fin_5521",
			AmountPaise:     12500000, // ₹1,25,000.00
			Currency:        "INR",
			Status:          models.StatusPending,
			ReferenceID:     "RTGS/R202609240092/EQUITY_DEP",
			Type:            models.TypeCredit,
			CreatedAt:       now.Add(-10 * time.Minute),
			FormattedAmount: models.FormatRupees(12500000),
		},
		{
			ID:              "a1029384-b56c-48de-9f12-000000000008",
			UserID:          "usr_fin_4412",
			AmountPaise:     49900, // ₹499.00
			Currency:        "INR",
			Status:          models.StatusFailed,
			ReferenceID:     "CARD/AUTH_8910283019",
			Type:            models.TypeDebit,
			CreatedAt:       now.Add(-6 * time.Hour),
			FormattedAmount: models.FormatRupees(49900),
		},
		{
			ID:              "a1029384-b56c-48de-9f12-000000000009",
			UserID:          "usr_fin_6634",
			AmountPaise:     850000, // ₹8,500.00
			Currency:        "INR",
			Status:          models.StatusSuccess,
			ReferenceID:     "UPI/428900192834/LOAN_EMI",
			Type:            models.TypeDebit,
			CreatedAt:       now.Add(-8 * time.Hour),
			FormattedAmount: models.FormatRupees(850000),
		},
		{
			ID:              "a1029384-b56c-48de-9f12-000000000010",
			UserID:          "usr_fin_2209",
			AmountPaise:     187525, // ₹1,875.25
			Currency:        "INR",
			Status:          models.StatusSuccess,
			ReferenceID:     "UPI/428938491029/GROCERY_EXP",
			Type:            models.TypeDebit,
			CreatedAt:       now.Add(-12 * time.Hour),
			FormattedAmount: models.FormatRupees(187525),
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
	var totalVolumePaise int64
	var pendingCount, successCount, failedCount int

	targetStatus := strings.ToUpper(strings.TrimSpace(statusFilter))

	for _, tx := range r.transactions {
		// Calculate summary statistics across all transactions in store
		totalVolumePaise += tx.AmountPaise
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
		TotalVolumePaise:       totalVolumePaise,
		TotalVolumeFormatted:   models.FormatRupees(totalVolumePaise),
		TotalCount:             totalCount,
		PendingCount:           pendingCount,
		SuccessCount:           successCount,
		FailedCount:            failedCount,
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
