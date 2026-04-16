package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/card-demo/project/src/app/report/service"
	"github.com/card-demo/project/src/app/report/dto"
	"github.com/card-demo/project/src/domain/report/repository"
	"github.com/go-chi/chi/v5"
)

// ReportHandler handles HTTP requests for Reports.
type ReportHandler struct {
	service *service.ReportService
}

// NewReportHandler creates a new ReportHandler.
func NewReportHandler(svc *service.ReportService) *ReportHandler {
	return &ReportHandler{service: svc}
}

// CreateReportRequestDTO defines the body for creating a report.
type CreateReportRequestDTO struct {
	Type   string            `json:"type"`
	Params map[string]string `json:"params"`
}

// CreateReport handles POST /reports
func (h *ReportHandler) CreateReport(w http.ResponseWriter, r *http.Request) {
	var req CreateReportRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Type == "" {
		http.Error(w, "Report type is required", http.StatusBadRequest)
		return
	}

	// Validate Date Params if present
	// (Basic validation for example)
	if startDate, ok := req.Params["startDate"]; ok {
		if _, err := time.Parse("2006-01-02", startDate); err != nil {
			http.Error(w, "Invalid startDate format (YYYY-MM-DD)", http.StatusBadRequest)
			return
		}
	}

	agg, err := h.service.CreateReport(r.Context(), req.Type, req.Params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := dto.ReportResponse{
		ID:        agg.ID,
		Type:      agg.Type,
		Status:    agg.Status,
		CreatedAt: agg.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(resp)
}

// GetReport handles GET /reports/{id}
func (h *ReportHandler) GetReport(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	agg, err := h.service.Get(r.Context(), id)
	if err != nil {
		if err.Error() == "not found" {
			http.Error(w, "Report not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := dto.ReportResponse{
		ID:        agg.ID,
		Type:      agg.Type,
		Status:    agg.Status,
		CreatedAt: agg.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// ListReports handles GET /reports
func (h *ReportHandler) ListReports(w http.ResponseWriter, r *http.Request) {
	// Parse Query Params
	filters := repository.ListFilters{
		Type:   r.URL.Query().Get("type"),
		Status: r.URL.Query().Get("status"),
	}

	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		filters.Page, _ = strconv.Atoi(pageStr)
	}
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		filters.Limit, _ = strconv.Atoi(limitStr)
	}

	// In a real app, default pagination logic would apply.

	list, err := h.service.List(r.Context(), filters)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := make([]dto.ReportResponse, len(list))
	for i, agg := range list {
		resp[i] = dto.ReportResponse{
			ID:        agg.ID,
			Type:      agg.Type,
			Status:    agg.Status,
			CreatedAt: agg.CreatedAt,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// DeleteReport handles DELETE /reports/{id}
func (h *ReportHandler) DeleteReport(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := h.service.Delete(r.Context(), id)
	if err != nil {
		if err.Error() == "not found" {
			http.Error(w, "Report not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
