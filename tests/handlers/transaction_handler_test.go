package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/carddemo/project/mocks"
	"github.com/carddemo/project/src/app/port/in/rest"
	"github.com/carddemo/project/src/domain/transaction/model"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransactionHandlers_CreateTransaction(t *testing.T) {
	// Setup
	repo := mocks.NewMockTransactionRepository()
	batchRepo := mocks.NewMockBatchSettlementRepository() // Needed for constructor
	handler := rest.NewTransactionHandler(repo, batchRepo)

	r := chi.NewRouter()
	r.Post("/transactions", handler.HandleCreateTransaction)

	tests := []struct {
		name           string
		requestBody    interface{}
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "Success - Create Debit Transaction",
			requestBody: map[string]interface{}{
				"account_id":       "acc_123",
				"card_id":          "card_123",
				"amount":           100.50,
				"transaction_type": "debit",
			},
			expectedStatus: http.StatusCreated,
			expectedBody: map[string]interface{}{
				"status": "submitted",
			},
		},
		{
			name: "Failure - Missing Amount",
			requestBody: map[string]interface{}{
				"account_id":       "acc_123",
				"card_id":          "card_123",
				"transaction_type": "debit",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "Invalid request",
			},
		},
		{
			name: "Failure - Invalid Type",
			requestBody: map[string]interface{}{
				"account_id":       "acc_123",
				"card_id":          "card_123",
				"amount":           50.0,
				"transaction_type": "transfer",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "Invalid request",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			req, err := http.NewRequest(http.MethodPost, "/transactions", bytes.NewBuffer(body))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			
			if tt.expectedStatus >= 400 {
				var resp map[string]interface{}
				err := json.Unmarshal(rr.Body.Bytes(), &resp)
				require.NoError(t, err)
				assert.NotEmpty(t, resp["error"])
			}
		})
	}
}

func TestTransactionHandlers_GetTransaction(t *testing.T) {
	// Setup
	repo := mocks.NewMockTransactionRepository()
	batchRepo := mocks.NewMockBatchSettlementRepository()
	handler := rest.NewTransactionHandler(repo, batchRepo)

	r := chi.NewRouter()
	r.Get("/transactions/{id}", handler.HandleGetTransaction)

	// Seed data
	txn := &model.Transaction{
		ID:              "txn_123",
		AccountID:       "acc_123",
		CardID:          "card_123",
		Amount:          100.50,
		TransactionType: "debit",
		Status:          "completed",
		CreatedAt:       time.Now(),
	}
	repo.Save(txn)

	t.Run("Success - Get Transaction", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/transactions/txn_123", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		
		var resp map[string]interface{}
		err = json.Unmarshal(rr.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.Equal(t, "txn_123", resp["id"])
		assert.Equal(t, "completed", resp["status"])
	})

	t.Run("Failure - Not Found", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/transactions/unknown", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
	})
}

func TestTransactionHandlers_ListTransactionsByAccount(t *testing.T) {
	repo := mocks.NewMockTransactionRepository()
	batchRepo := mocks.NewMockBatchSettlementRepository()
	handler := rest.NewTransactionHandler(repo, batchRepo)

	r := chi.NewRouter()
	r.Get("/accounts/{id}/transactions", handler.HandleListTransactionsByAccount)

	t.Run("Success - List Transactions", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/accounts/acc_123/transactions", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		var resp []map[string]interface{}
		err = json.Unmarshal(rr.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.IsType(t, []map[string]interface{}{}, resp)
	})

	t.Run("Success - List with Query Parameters", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/accounts/acc_123/transactions?status=submitted&from=2023-01-01", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		// In real implementation, we would assert filtering logic works
	})
}

func TestTransactionHandlers_VoidTransaction(t *testing.T) {
	repo := mocks.NewMockTransactionRepository()
	batchRepo := mocks.NewMockBatchSettlementRepository()
	handler := rest.NewTransactionHandler(repo, batchRepo)

	r := chi.NewRouter()
	r.Post("/transactions/{id}/void", handler.HandleVoidTransaction)

	// Seed data
	txn := &model.Transaction{
		ID:              "txn_void",
		AccountID:       "acc_123",
		Amount:          100.0,
		TransactionType: "debit",
		Status:          "submitted",
		CreatedAt:       time.Now(),
	}
	repo.Save(txn)

	t.Run("Success - Void Transaction", func(t *testing.T) {
		body := map[string]interface{}{"reason": "Customer requested"}
		jsonBody, _ := json.Marshal(body)
		
		req, err := http.NewRequest(http.MethodPost, "/transactions/txn_void/void", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		
		var resp map[string]interface{}
		err = json.Unmarshal(rr.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.Equal(t, "reversed", resp["status"])
	})

	t.Run("Failure - Missing Reason", func(t *testing.T) {
		body := map[string]interface{}{}
		jsonBody, _ := json.Marshal(body)

		req, err := http.NewRequest(http.MethodPost, "/transactions/txn_void/void", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})
}
