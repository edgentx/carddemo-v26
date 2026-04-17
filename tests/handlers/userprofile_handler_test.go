package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/carddemo/project/src/app/port/in/rest"
	"github.com/carddemo/project/src/domain/userprofile/model"
	"github.com/carddemo/project/mocks"
	"github.com/go-chi/chi/v5"
)

func setupProfileHandlerTest(repo *mocks.MockUserProfileRepository) *chi.Mux {
	// Note: In a real scenario, this would wire up a ProfileService.
	// For this test suite, we assume the wiring happens in rest.RegisterRoutes or similar.
	// We verify the endpoints exist and parse inputs correctly.
	handler := rest.NewProfileHandler(repo) // Assuming constructor exists

	r := chi.NewRouter()
	r.Get("/accounts/{id}/profile", handler.HandleGetProfile) // Assuming direct registration or helper
	r.Put("/accounts/{id}/profile", handler.HandleUpdateProfile)
	return r
}

func TestUserProfileEndpoints_HandleGetProfile(t *testing.T) {
	t.Run("Success: Returns 200 with profile data", func(t *testing.T) {
		// Mock Setup
		repo := mocks.NewMockUserProfileRepository()
		profile, _ := model.NewUserProfile("profile-01", "John", "Doe", "john@example.com")
		repo.Save(profile)

		r := setupProfileHandlerTest(repo)

		req := httptest.NewRequest("GET", "/accounts/profile-01/profile", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		if resp["email"] != "john@example.com" {
			t.Errorf("Expected email john@example.com, got %v", resp["email"])
		}
	})

	t.Run("Failure: Returns 404 for non-existent profile", func(t *testing.T) {
		repo := mocks.NewMockUserProfileRepository()
		r := setupProfileHandlerTest(repo)

		req := httptest.NewRequest("GET", "/accounts/missing/profile", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", w.Code)
		}
	})
}

func TestUserProfileEndpoints_HandleUpdateProfile(t *testing.T) {
	t.Run("Success: Returns 200 with updated data", func(t *testing.T) {
		repo := mocks.NewMockUserProfileRepository()
		profile, _ := model.NewUserProfile("profile-02", "Jane", "Doe", "jane@example.com")
		repo.Save(profile)

		r := setupProfileHandlerTest(repo)

		body := map[string]string{"first_name": "Jane", "last_name": "Smith", "email": "jane.smith@example.com"}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest("PUT", "/accounts/profile-02/profile", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		if resp["last_name"] != "Smith" {
			t.Errorf("Expected last_name Smith, got %v", resp["last_name"])
		}
	})
}
