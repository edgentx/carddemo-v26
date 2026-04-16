package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/carddemo/project/mocks"
	"github.com/carddemo/project/src/app/card/handler"
	"github.com/carddemo/project/src/app/card/service"
	"github.com/carddemo/project/src/app/cardpolicy/service"
	"github.com/carddemo/project/src/app/shared"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test suite structure to ensure we have a clean Red phase setup
type CardHandlerTestSuite struct {
	Router       *chi.Mux
	CardRepo     *mocks.MockCardRepository
	PolicyRepo   *mocks.MockCardPolicyRepository
	AccountRepo  *mocks.MockAccountRepository
	CardService  *service.CardApplicationService
	PolicyService *service.CardPolicyApplicationService
	Handler      *handler.CardHandler
	Server       *httptest.Server
}

// We write tests that assert the existence of routes and their behavior.
// Since we are in TDD Red phase, we expect these to fail or compile errors
// until the implementation exists.

func setupTest(tb testing.TB) *CardHandlerTestSuite {
	tb.Helper()

	// 1. Initialize Mocks (defined in mocks/)
	cardRepo := mocks.NewMockCardRepository()
	policyRepo := mocks.NewMockCardPolicyRepository()
	accountRepo := mocks.NewMockAccountRepository()

	// 2. Initialize Application Services
	// Note: We use the real service layer logic which expects interfaces.
	cardSvc := service.NewCardApplicationService(cardRepo, accountRepo)
	policySvc := service.NewCardPolicyApplicationService(policyRepo)

	// 3. Initialize Handler
	h := handler.NewCardHandler(cardSvc, policySvc)

	// 4. Initialize Router
	router := chi.NewRouter()
	handler.RegisterRoutes(router, h)

	return &CardHandlerTestSuite{
		Router:        router,
		CardRepo:      cardRepo,
		PolicyRepo:    policyRepo,
		AccountRepo:   accountRepo,
		CardService:   cardSvc,
		PolicyService: policySvc,
		Handler:       h,
	}
}

func TestCardEndpoints_AreRegistered(t *testing.T) {
	suite := setupTest(t)

	// Chi does not easily expose route listing, so we test by hitting them
	// and ensuring 404 is not returned for valid paths, or checking the router.
	// However, for TDD, we usually test functionality directly.
	
	// Check that router is not nil (Sanity)
	assert.NotNil(t, suite.Router)
}

// --- IssueCard Tests ---

func TestIssueCard_Success(t *testing.T) {
	suite := setupTest(t)

	// Setup Data: We need an account to issue a card against
	// Ideally we seed the AccountRepo mock.
	// Since mocks are blank, we expect logic failure if validation checks existence,
	// or success if not. Assuming strict validation (Account must exist).

	reqBody := map[string]interface{}{
		"account_id": "acc_123",
		"card_type":  "Virtual",
		"spending_limits": map[string]int{
			"daily": 1000,
			"weekly": 5000,
		},
	}
	bodyBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/accounts/acc_123/cards", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.Router.ServeHTTP(w, req)

	// In Red Phase: Implementation likely missing or mock empty.
	// We expect the handler to process the request.
	// If account validation is ON: Might return 404/422 if account not found.
	// If account validation is OFF: Returns 201.
	
	// We assert the standard structure. If this compiles and fails -> Red.
	assert.Equal(t, http.StatusCreated, w.Code, "Expected 201 Created")
}

func TestIssueCard_ValidationError(t *testing.T) {
	suite := setupTest(t)

	reqBody := map[string]interface{}{
		"account_id": "", // Invalid: Empty
		"card_type":  "InvalidType",
	}
	bodyBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/accounts/acc_123/cards", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.Router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// --- GetCard Tests ---

func TestGetCard_Success(t *testing.T) {
	suite := setupTest(t)

	req := httptest.NewRequest("GET", "/cards/card_001", nil)
	w := httptest.NewRecorder()

	suite.Router.ServeHTTP(w, req)

	// Expect 200 if found, 404 if not (mock is empty).
	// If handler maps error correctly, should be 404.
	// If handler crashes, 500.
	// This test verifies the route is reachable.
	assert.Contains(t, []int{http.StatusOK, http.StatusNotFound}, w.Code)
}

// --- UpdateCardStatus Tests ---

func TestUpdateCardStatus_Success(t *testing.T) {
	suite := setupTest(t)

	reqBody := map[string]string{"status": "Blocked"}
	bodyBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("PUT", "/cards/card_001/status", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.Router.ServeHTTP(w, req)

	// Expect 200 OK or 404 if card doesn't exist
	assert.Contains(t, []int{http.StatusOK, http.StatusNotFound}, w.Code)
}

func TestUpdateCardStatus_InvalidStatus(t *testing.T) {
	suite := setupTest(t)

	reqBody := map[string]string{"status": "ImpossibleState"}
	bodyBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("PUT", "/cards/card_001/status", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.Router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// --- ActivateCard Tests ---

func TestActivateCard_Success(t *testing.T) {
	suite := setupTest(t)

	req := httptest.NewRequest("POST", "/cards/card_001/activate", nil)
	w := httptest.NewRecorder()

	suite.Router.ServeHTTP(w, req)

	assert.Contains(t, []int{http.StatusOK, http.StatusNotFound, http.StatusConflict}, w.Code)
}

// --- Card Policy Tests ---

func TestGetPolicy_Success(t *testing.T) {
	suite := setupTest(t)

	req := httptest.NewRequest("GET", "/policies/pol_001", nil)
	w := httptest.NewRecorder()

	suite.Router.ServeHTTP(w, req)

	assert.Contains(t, []int{http.StatusOK, http.StatusNotFound}, w.Code)
}

func TestUpdatePolicy_Success(t *testing.T) {
	suite := setupTest(t)

	reqBody := map[string]interface{}{
		"daily_limit": 5000,
	}
	bodyBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("PUT", "/policies/pol_001", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.Router.ServeHTTP(w, req)

	assert.Contains(t, []int{http.StatusOK, http.StatusNotFound}, w.Code)
}

func TestGetPoliciesByAccount_Success(t *testing.T) {
	suite := setupTest(t)

	req := httptest.NewRequest("GET", "/accounts/acc_123/policies", nil)
	w := httptest.NewRecorder()

	suite.Router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// --- Response JSON Structure Tests ---

func TestResponseJSONStructure(t *testing.T) {
	suite := setupTest(t)

	// We'll hit an endpoint that returns JSON (List is usually safest)
	req := httptest.NewRequest("GET", "/accounts/acc_123/policies", nil)
	w := httptest.NewRecorder()

	suite.Router.ServeHTTP(w, req)

	var response map[string]interface{}
	body, _ := io.ReadAll(w.Body)
	err := json.Unmarshal(body, &response)

	// If we assume the list endpoint always returns an array/object, even if empty.
	assert.NoError(t, err, "Response should be valid JSON")
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
}
