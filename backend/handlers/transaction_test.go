package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/RR-Sahoo/fintech-starter-kit/backend/handlers"
	"github.com/RR-Sahoo/fintech-starter-kit/backend/models"
	"github.com/RR-Sahoo/fintech-starter-kit/backend/repository"
)

func setupTestRouter() (*gin.Engine, *repository.MemoryTransactionRepository) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	repo := repository.NewMemoryTransactionRepository()
	h := handlers.NewTransactionHandler(repo)

	v1 := r.Group("/api/v1")
	{
		v1.GET("/transactions", h.GetTransactions)
		v1.POST("/transactions/reconcile", h.ReconcileTransaction)
	}

	return r, repo
}

func TestGetTransactions(t *testing.T) {
	router, _ := setupTestRouter()

	t.Run("Get all transactions with summary", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/transactions", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", w.Code)
		}

		var response struct {
			Transactions []models.Transaction         `json:"transactions"`
			Summary      repository.TransactionSummary `json:"summary"`
		}
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if len(response.Transactions) == 0 {
			t.Errorf("expected non-empty transactions")
		}
		if response.Summary.TotalCount != len(response.Transactions) {
			t.Errorf("expected summary count %d, got %d", len(response.Transactions), response.Summary.TotalCount)
		}
		if response.Summary.PendingCount == 0 {
			t.Errorf("expected at least one pending transaction")
		}
	})

	t.Run("Filter by status PENDING", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/transactions?status=PENDING", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", w.Code)
		}

		var response struct {
			Transactions []models.Transaction `json:"transactions"`
		}
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		for _, tx := range response.Transactions {
			if tx.Status != models.StatusPending {
				t.Errorf("expected status PENDING, got %s", tx.Status)
			}
		}
	})

	t.Run("Invalid status filter", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/transactions?status=UNKNOWN", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected status 400, got %d", w.Code)
		}
	})

	t.Run("Limit query parameter", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/transactions?limit=2", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", w.Code)
		}

		var response struct {
			Transactions []models.Transaction `json:"transactions"`
		}
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if len(response.Transactions) != 2 {
			t.Errorf("expected 2 transactions, got %d", len(response.Transactions))
		}
	})
}

func TestReconcileTransaction(t *testing.T) {
	router, _ := setupTestRouter()

	t.Run("Reconcile existing pending transaction successfully", func(t *testing.T) {
		// id 3 is pending in mock store
		payload := []byte(`{"transaction_id": "a1029384-b56c-48de-9f12-000000000003"}`)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/transactions/reconcile", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d: %s", w.Code, w.Body.String())
		}

		var response struct {
			Message     string             `json:"message"`
			Transaction models.Transaction `json:"transaction"`
		}
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if response.Transaction.Status != models.StatusSuccess {
			t.Errorf("expected status SUCCESS, got %s", response.Transaction.Status)
		}

		// Second reconcile attempt on same transaction should now return 409 Conflict
		w2 := httptest.NewRecorder()
		req2, _ := http.NewRequest(http.MethodPost, "/api/v1/transactions/reconcile", bytes.NewBuffer(payload))
		req2.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w2, req2)
		if w2.Code != http.StatusConflict {
			t.Errorf("expected status 409 on already reconciled transaction, got %d", w2.Code)
		}
	})

	t.Run("Reconcile non-existent transaction returns 404", func(t *testing.T) {
		payload := []byte(`{"transaction_id": "non-existent-id"}`)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/transactions/reconcile", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Fatalf("expected status 404, got %d", w.Code)
		}
	})

	t.Run("Missing transaction_id returns 400", func(t *testing.T) {
		payload := []byte(`{}`)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/transactions/reconcile", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected status 400, got %d", w.Code)
		}
	})
}
