package handler

import (
	"encoding/json"
	"net/http"

	"github.com/card-demo/project/src/app/exportjob/service"
	"github.com/card-demo/project/src/app/exportjob/dto"
	"github.com/card-demo/project/src/app/shared/storage"
	"github.com/go-chi/chi/v5"
)

// ExportHandler handles HTTP requests for Export Jobs.
type ExportHandler struct {
	service      *service.ExportService
	storageClient storage.StorageClient
}

// NewExportHandler creates a new ExportHandler.
// storageClient is used for downloading files.
func NewExportHandler(svc *service.ExportService, storageClient storage.StorageClient) *ExportHandler {
	return &ExportHandler{
		service:      svc,
		storageClient: storageClient,
	}
}

// CreateExportRequestDTO defines the body for creating an export.
type CreateExportRequestDTO struct {
	ReportID string `json:"report_id"`
}

// CreateExport handles POST /exports
func (h *ExportHandler) CreateExport(w http.ResponseWriter, r *http.Request) {
	var req CreateExportRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.ReportID == "" {
		http.Error(w, "report_id is required", http.StatusBadRequest)
		return
	}

	// Call Service
	agg, err := h.service.CreateExport(r.Context(), req.ReportID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Map to Response DTO
	resp := dto.ExportJobResponse{
		ID:        agg.ID,
		ReportID:  agg.ReportID,
		Status:    agg.Status,
		CreatedAt: agg.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(resp)
}

// GetExport handles GET /exports/{id}
func (h *ExportHandler) GetExport(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	
	agg, err := h.service.Get(r.Context(), id)
	if err != nil {
		// Assuming domain.ErrNotFound is handled appropriately or mapped here
		if err.Error() == "not found" { // Simplistic check, rely on domain errors in prod
			http.Error(w, "Export job not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := dto.ExportJobResponse{
		ID:        agg.ID,
		ReportID:  agg.ReportID,
		Status:    agg.Status,
		FileKey:   agg.FileKey,
		CreatedAt: agg.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// DownloadExport handles GET /exports/{id}/download
func (h *ExportHandler) DownloadExport(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	// 1. Get Job Status
	agg, err := h.service.Get(r.Context(), id)
	if err != nil {
		http.Error(w, "Export job not found", http.StatusNotFound)
		return
	}

	// 2. Check if ready
	if agg.Status != "completed" {
		http.Error(w, "Export job not ready", http.StatusAccepted)
		return
	}

	// 3. Stream from Storage
	data, contentType, err := h.storageClient.GetFile(agg.FileKey)
	if err != nil {
		http.Error(w, "Failed to retrieve file", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Disposition", "attachment; filename="+id+".csv")
	w.Write(data)
}
