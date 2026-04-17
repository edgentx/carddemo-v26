package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/carddemo/project/src/app/account/dto"
	"github.com/carddemo/project/src/app/account/service"
	"github.com/carddemo/project/src/app/port/in/rest"
	"github.com/carddemo/project/src/domain/account/command"
	"github.com/carddemo/project/src/domain/account/model"
	accountRepo "github.com/carddemo/project/src/domain/account/repository"
	"github.com/carddemo/project/src/domain/shared"
	"github.com/carddemo/project/mocks"
	"github.com/go-chi/chi/v5"
)

// setupHandlerTest initializes the router with mock dependencies.
func setupHandlerTest(repo *mocks.MockAccountRepository) *chi.Mux {
	// We use the real Service layer wired with the Mock Repository to satisfy
	// the integration requirement while mocking external dependencies (DB).
	acctService := service.NewAccountService(repo)

	handler := rest.NewAccountHandler(acctService)

	r := chi.NewRouter()
	rest.RegisterRoutes(r, handler)
	return r
}

func TestAccountEndpoints_HandleOpenAccount(t *testing.T) {
	t.Run("Success: Returns 201 and Location header", func(t *testing.T) {
		repo := mocks.NewMockAccountRepository()
		r := setupHandlerTest(repo)

		body := map[string]interface{}{
			"user_profile_id": "user-123",
			"status":          "active",
			"account_type":    "savings",
		}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest("POST", "/accounts", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("Expected status 201, got %d", w.Code)
		}

		loc := w.Header().Get("Location")
		if loc == "" {
			t.Error("Expected Location header to be set")
		}

		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		if resp["id"] == nil {
			t.Error("Expected ID in response")
		}
	})

	t.Run("Validation: Returns 400 for invalid JSON", func(t *testing.T) {
		repo := mocks.NewMockAccountRepository()
		r := setupHandlerTest(repo)

		req := httptest.NewRequest("POST", "/accounts", bytes.NewBuffer([]byte("{invalid")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400 for bad JSON, got %d", w.Code)
		}
	})
}

func TestAccountEndpoints_HandleGetAccount(t *testing.T) {
	t.Run("Success: Returns 200 with account details", func(t *testing.T) {
		repo := mocks.NewMockAccountRepository()
		r := setupHandlerTest(repo)

		// Setup: Create an account via the mock repo directly to simulate DB state
		// This satisfies the "Aggregate State" requirement without needing a POST
		// to pre-seed, isolating the GET logic.
		aggregate, _ := model.NewAccount("acct-01", "user-01", "pending", "checking")
		// Apply Timestamps
		now := time.Now()
		aggregate.CreatedAt = now
		aggregate.UpdatedAt = now
		repo.Save(aggregate)

		req := httptest.NewRequest("GET", "/accounts/acct-01", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)

		if resp["id"] != "acct-01" {
			t.Errorf("Expected ID acct-01, got %v", resp["id"])
		}
		if resp["created_at"] == nil {
			t.Error("Expected created_at timestamp in response")
		}
	})

	t.Run("Failure: Returns 404 when account not found", func(t *testing.T) {
		repo := mocks.NewMockAccountRepository()
		r := setupHandlerTest(repo)

		req := httptest.NewRequest("GET", "/accounts/non-existent", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", w.Code)
		}
	})
}

func TestAccountEndpoints_HandleUpdateAccountStatus(t *testing.T) {
	t.Run("Success: Returns 200 on status update", func(t *testing.T) {
		repo := mocks.NewMockAccountRepository()
		r := setupHandlerTest(repo)

		aggregate, _ := model.NewAccount("acct-02", "user-02", "pending", "checking")
		repo.Save(aggregate)

		body := map[string]string{"new_status": "active", "reason": "KYC passed"}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest("PUT", "/accounts/acct-02/status", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}

		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		if resp["status"] != "active" {
			t.Errorf("Expected status active, got %v", resp["status"])
		}
	})

	t.Run("Failure: Returns 404 if account missing", func(t *testing.T) {
		repo := mocks.NewMockAccountRepository()
		r := setupHandlerTest(repo)

		body := map[string]string{"new_status": "active"}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest("PUT", "/accounts/missing/status", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", w.Code)
		}
	})
}

func TestAccountEndpoints_HandleDeleteAccount(t *testing.T) {
	t.Run("Success: Returns 204 No Content", func(t *testing.T) {
		repo := mocks.NewMockAccountRepository()
		r := setupHandlerTest(repo)

		aggregate, _ := model.NewAccount("acct-del", "user-del", "active", "checking")
		repo.Save(aggregate)

		req := httptest.NewRequest("DELETE", "/accounts/acct-del", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusNoContent {
			t.Errorf("Expected status 204, got %d", w.Code)
		}
		if w.Body.Len() != 0 {
			t.Errorf("Expected empty body, got %s", w.Body.String())
		}
	})
}
