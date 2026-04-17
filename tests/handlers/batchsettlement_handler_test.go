package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/carddemo/project/mocks"
	"github.com/carddemo/project/src/app/port/in/rest"
	"github.com/carddemo/project/src/domain/batchsettlement/model"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBatchSettlementHandlers_CreateBatch(t *testing.T) {
	txnRepo := mocks.NewMockTransactionRepository()
	batchRepo := mocks.NewMockBatchSettlementRepository()
	handler := rest.NewTransactionHandler(txnRepo, batchRepo)

	r := chi.NewRouter()
	r.Post("/settlements", handler.HandleCreateBatchSettlement)

	t.Run("Success - Create Batch", func(t *testing.T) {
		body := map[string]interface{}{
			"name":        "Daily Closeout",
			"description": "End of day processing",
		}
		jsonBody, _ := json.Marshal(body)

		req, err := http.NewRequest(http.MethodPost, "/settlements", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)

		var resp map[string]interface{}
		err = json.Unmarshal(rr.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.Equal(t, "Daily Closeout", resp["name"])
		assert.Equal(t, "open", resp["status"])
	})

	t.Run("Failure - Missing Name", func(t *testing.T) {
		body := map[string]interface{}{"description": "test"}
		jsonBody, _ := json.Marshal(body)

		req, err := http.NewRequest(http.MethodPost, "/settlements", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})
}

func TestBatchSettlementHandlers_GetBatch(t *testing.T) {
	txnRepo := mocks.NewMockTransactionRepository()
	batchRepo := mocks.NewMockBatchSettlementRepository()
	handler := rest.NewTransactionHandler(txnRepo, batchRepo)

	r := chi.NewRouter()
	r.Get("/settlements/{id}", handler.HandleGetBatchSettlement)

	t.Run("Success - Get Batch", func(t *testing.T) {
		// Seed data
		batch := &model.BatchSettlement{
			ID:          "batch_1",
			Name:        "Batch 1",
			Description: "Desc",
			Status:      "reconciled",
			CreatedAt:   time.Now(),
		}
		batchRepo.Save(batch)

		req, err := http.NewRequest(http.MethodGet, "/settlements/batch_1", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		var resp map[string]interface{}
		err = json.Unmarshal(rr.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.Equal(t, "batch_1", resp["id"])
	})

	t.Run("Failure - Not Found", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/settlements/unknown", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
	})
}

func TestBatchSettlementHandlers_ListBatches(t *testing.T) {
	txnRepo := mocks.NewMockTransactionRepository()
	batchRepo := mocks.NewMockBatchSettlementRepository()
	handler := rest.NewTransactionHandler(txnRepo, batchRepo)

	r := chi.NewRouter()
	r.Get("/settlements", handler.HandleListBatchSettlements)

	t.Run("Success - List Batches", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/settlements", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		var resp []map[string]interface{}
		err = json.Unmarshal(rr.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.IsType(t, []map[string]interface{}{}, resp)
	})
}
