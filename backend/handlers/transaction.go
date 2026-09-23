package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/RR-Sahoo/fintech-starter-kit/backend/models"
	"github.com/RR-Sahoo/fintech-starter-kit/backend/repository"
)

// TransactionHandler handles HTTP requests for transactions.
type TransactionHandler struct {
	repo repository.TransactionRepository
}

// NewTransactionHandler creates a new TransactionHandler instance.
func NewTransactionHandler(repo repository.TransactionRepository) *TransactionHandler {
	return &TransactionHandler{repo: repo}
}

// ReconcileRequest represents the request body for reconciling a transaction.
type ReconcileRequest struct {
	TransactionID string `json:"transaction_id" binding:"required"`
}

// GetTransactions handles GET /api/v1/transactions
// Query parameters:
// - status: filter by SUCCESS, PENDING, or FAILED (case-insensitive)
// - limit: max number of transactions to return (e.g. 10, 20, 50)
func (h *TransactionHandler) GetTransactions(c *gin.Context) {
	statusFilter := strings.ToUpper(strings.TrimSpace(c.Query("status")))
	if statusFilter != "" {
		if statusFilter != models.StatusSuccess &&
			statusFilter != models.StatusPending &&
			statusFilter != models.StatusFailed {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid status filter. Must be one of: SUCCESS, PENDING, FAILED",
			})
			return
		}
	}

	limitStr := strings.TrimSpace(c.Query("limit"))
	limit := 0
	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err != nil || parsedLimit <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Limit must be a positive integer",
			})
			return
		}
		limit = parsedLimit
	}

	transactions, summary := h.repo.GetAll(statusFilter, limit)

	c.JSON(http.StatusOK, gin.H{
		"transactions": transactions,
		"summary":      summary,
	})
}

// ReconcileTransaction handles POST /api/v1/transactions/reconcile
// Accepts { "transaction_id": "..." } and updates status from PENDING to SUCCESS.
func (h *TransactionHandler) ReconcileTransaction(c *gin.Context) {
	var req ReconcileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing or invalid 'transaction_id' in request body",
		})
		return
	}

	transactionID := strings.TrimSpace(req.TransactionID)
	if transactionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "'transaction_id' cannot be empty",
		})
		return
	}

	updatedTx, err := h.repo.Reconcile(transactionID)
	if err != nil {
		if errors.Is(err, repository.ErrTransactionNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Transaction not found with the provided ID",
			})
			return
		}
		if errors.Is(err, repository.ErrTransactionNotPending) {
			c.JSON(http.StatusConflict, gin.H{
				"error": "Only transactions with status PENDING can be reconciled",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to reconcile transaction",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Transaction successfully reconciled",
		"transaction": updatedTx,
	})
}
