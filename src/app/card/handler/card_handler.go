package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/carddemo/project/src/app/card/dto"
	"github.com/carddemo/project/src/app/card/service"
	policyService "github.com/carddemo/project/src/app/cardpolicy/service"
	"github.com/go-chi/chi/v5"
)

// CardHandler handles HTTP requests for Card and CardPolicy aggregates.
type CardHandler struct {
	cardSvc   *service.CardApplicationService
	policySvc *policyService.CardPolicyApplicationService
}

// NewCardHandler creates a new CardHandler.
func NewCardHandler(cardSvc *service.CardApplicationService, policySvc *policyService.CardPolicyApplicationService) *CardHandler {
	return &CardHandler{
		cardSvc:   cardSvc,
		policySvc: policySvc,
	}
}

// RegisterRoutes registers the card and policy routes on the router.
func RegisterRoutes(r chi.Router, h *CardHandler) {
	r.Route("/accounts", func(r chi.Router) {
		r.Route("/{accountID}", func(r chi.Router) {
			r.Post("/cards", h.IssueCard)
			r.Get("/policies", h.GetPoliciesByAccount)
		})
	})

	r.Route("/cards", func(r chi.Router) {
		r.Route("/{cardID}", func(r chi.Router) {
			r.Get("/", h.GetCard)
			r.Put("/status", h.UpdateCardStatus)
			r.Post("/activate", h.ActivateCard)
		})
	})

	r.Route("/policies", func(r chi.Router) {
		r.Route("/{policyID}", func(r chi.Router) {
			r.Get("/", h.GetPolicy)
			r.Put("/", h.UpdatePolicy)
		})
	})
}

// --- Card Handlers ---

// IssueCard handles POST /accounts/{id}/cards
func (h *CardHandler) IssueCard(w http.ResponseWriter, r *http.Request) {
	accountID := chi.URLParam(r, "accountID")
	if accountID == "" {
		http.Error(w, "account_id is required", http.StatusBadRequest)
		return
	}

	var req dto.IssueCardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Basic validation
	if req.CardType == "" {
		http.Error(w, "card_type is required", http.StatusBadRequest)
		return
	}

	cardID, err := h.cardSvc.IssueCard(r.Context(), accountID, req.CardType, req.SpendingLimits)
	if err != nil {
		// In a real app, check error types for specific status codes
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": cardID})
}

// GetCard handles GET /cards/{id}
func (h *CardHandler) GetCard(w http.ResponseWriter, r *http.Request) {
	cardID := chi.URLParam(r, "cardID")
	if cardID == "" {
		http.Error(w, "card_id is required", http.StatusBadRequest)
		return
	}

	card, err := h.cardSvc.GetCard(r.Context(), cardID)
	if err != nil {
		if errors.Is(err, service.ErrCardNotFound) {
			http.Error(w, "card not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(card)
}

// UpdateCardStatus handles PUT /cards/{id}/status
func (h *CardHandler) UpdateCardStatus(w http.ResponseWriter, r *http.Request) {
	cardID := chi.URLParam(r, "cardID")
	if cardID == "" {
		http.Error(w, "card_id is required", http.StatusBadRequest)
		return
	}

	var req dto.UpdateCardStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Status == "" {
		http.Error(w, "status is required", http.StatusBadRequest)
		return
	}

	err := h.cardSvc.UpdateCardStatus(r.Context(), cardID, req.Status)
	if err != nil {
		if errors.Is(err, service.ErrCardNotFound) {
			http.Error(w, "card not found", http.StatusNotFound)
			return
		}
		// Assuming invalid status maps to Bad Request or Internal depending on implementation
		if err.Error() == "invalid card status" {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// ActivateCard handles POST /cards/{id}/activate
func (h *CardHandler) ActivateCard(w http.ResponseWriter, r *http.Request) {
	cardID := chi.URLParam(r, "cardID")
	if cardID == "" {
		http.Error(w, "card_id is required", http.StatusBadRequest)
		return
	}

	err := h.cardSvc.ActivateCard(r.Context(), cardID)
	if err != nil {
		if errors.Is(err, service.ErrCardNotFound) {
			http.Error(w, "card not found", http.StatusNotFound)
			return
		}
		if errors.Is(err, service.ErrInvalidStateTransition) {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// --- Policy Handlers ---

// GetPolicy handles GET /policies/{id}
func (h *CardHandler) GetPolicy(w http.ResponseWriter, r *http.Request) {
	policyID := chi.URLParam(r, "policyID")
	if policyID == "" {
		http.Error(w, "policy_id is required", http.StatusBadRequest)
		return
	}

	policy, err := h.policySvc.GetPolicy(r.Context(), policyID)
	if err != nil {
		if errors.Is(err, policyService.ErrPolicyNotFound) {
			http.Error(w, "policy not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(policy)
}

// UpdatePolicy handles PUT /policies/{id}
func (h *CardHandler) UpdatePolicy(w http.ResponseWriter, r *http.Request) {
	policyID := chi.URLParam(r, "policyID")
	if policyID == "" {
		http.Error(w, "policy_id is required", http.StatusBadRequest)
		return
	}

	var req dto.UpdatePolicyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err := h.policySvc.UpdatePolicyLimits(r.Context(), policyID, req.DailyLimit)
	if err != nil {
		if errors.Is(err, policyService.ErrPolicyNotFound) {
			http.Error(w, "policy not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// GetPoliciesByAccount handles GET /accounts/{id}/policies
func (h *CardHandler) GetPoliciesByAccount(w http.ResponseWriter, r *http.Request) {
	accountID := chi.URLParam(r, "accountID")
	if accountID == "" {
		http.Error(w, "account_id is required", http.StatusBadRequest)
		return
	}

	policies, err := h.policySvc.ListPoliciesByAccount(r.Context(), accountID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"policies": policies})
}
