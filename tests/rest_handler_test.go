package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/card-demo/project/src/app/report/dto"
	"github.com/card-demo/project/src/app/exportjob/dto"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock Repositories and Clients
// (Assuming MockExportJobRepository and MockReportRepository are generated or defined in mocks/)
// We define inline here for test file portability if mocks are not in separate files for this specific test run.

// --- Report Handler Tests ---

func TestReportHandler_CreateReport(t *testing.T) {
	mockRepo := new(MockReportRepository)
	handler := NewReportHandler(mockRepo)

	tests := []struct {
		name           string
		requestBody    interface{}
		mockReturn     *domain.ReportAggregate
		mockError      error
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Success - Create Report",
			requestBody: map[string]interface{}{
				"type": "transaction_summary",
				"params": map[string]string{
					"startDate": "2023-01-01",
					"endDate":   "2023-01-31",
				},
			},
			mockReturn: &domain.ReportAggregate{
				ID:     "report-123",
				Type:   "transaction_summary",
				Status: "pending",
			},
			mockError:      nil,
			expectedStatus: http.StatusAccepted,
			expectedBody:   `"id":"report-123"`,
		},
		{
			name:           "Failure - Invalid JSON",
			requestBody:    "invalid json",
			mockReturn:     nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid request body",
		},
		{
			name: "Failure - Missing Type",
			requestBody: map[string]interface{}{
				"params": map[string]string{},
			},
			mockReturn:     nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Report type is required",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var body io.Reader
			if str, ok := tc.requestBody.(string); ok {
				body = strings.NewReader(str)
			} else {
				b, _ := json.Marshal(tc.requestBody)
				body = bytes.NewReader(b)
			}

			req := httptest.NewRequest("POST", "/reports", body)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			if tc.mockReturn != nil || tc.mockError != nil {
				mockRepo.On("Save", mock.Anything, mock.Anything).Return(tc.mockError)
			}

			handler.CreateReport(w, req)

			resp := w.Result()
			respBody, _ := io.ReadAll(resp.Body)

			assert.Equal(t, tc.expectedStatus, resp.StatusCode)
			if tc.expectedBody != "" {
				assert.Contains(t, string(respBody), tc.expectedBody)
			}
		})
	}
}

func TestReportHandler_GetReport(t *testing.T) {
	mockRepo := new(MockReportRepository)
	handler := NewReportHandler(mockRepo)

	tests := []struct {
		name           string
		reportID       string
		mockReturn     *domain.ReportAggregate
		mockError      error
		expectedStatus int
	}{
		{
			name:     "Success",
			reportID: "report-123",
			mockReturn: &domain.ReportAggregate{
				ID:     "report-123",
				Status: "completed",
				Type:   "transaction_summary",
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Not Found",
			reportID:       "unknown",
			mockReturn:     nil,
			mockError:      domain.ErrNotFound,
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/reports/"+tc.reportID, nil)
			w := httptest.NewRecorder()

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tc.reportID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			mockRepo.On("Get", mock.Anything, tc.reportID).Return(tc.mockReturn, tc.mockError)

			handler.GetReport(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)
		})
	}
}

func TestReportHandler_ListReports(t *testing.T) {
	mockRepo := new(MockReportRepository)
	handler := NewReportHandler(mockRepo)

	req := httptest.NewRequest("GET", "/reports?type=transaction_summary&status=completed&page=1&limit=10", nil)
	w := httptest.NewRecorder()

	mockRepo.On("List", mock.Anything, mock.AnythingOfType("repository.ListFilters")).Return([]domain.ReportAggregate{{ID: "r1"}}, nil)

	handler.ListReports(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestReportHandler_DeleteReport(t *testing.T) {
	mockRepo := new(MockReportRepository)
	handler := NewReportHandler(mockRepo)

	req := httptest.NewRequest("DELETE", "/reports/r1", nil)
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "r1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	mockRepo.On("Delete", mock.Anything, "r1").Return(nil)

	handler.DeleteReport(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

// --- Export Job Handler Tests ---

func TestExportHandler_CreateExport(t *testing.T) {
	mockRepo := new(MockExportJobRepository)
	mockTemporal := new(MockTemporalClient) // Defined in mocks/
	handler := NewExportHandler(mockRepo, mockTemporal)

	validBody := map[string]interface{}{"report_id": "r1"}
	body, _ := json.Marshal(validBody)

	req := httptest.NewRequest("POST", "/exports", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mockTemporal.On("ExecuteWorkflow", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)
	mockRepo.On("Save", mock.Anything, mock.Anything).Return(nil) // Assuming initial save happens in service

	handler.CreateExport(w, req)

	assert.Equal(t, http.StatusAccepted, w.Code)
}

func TestExportHandler_DownloadExport(t *testing.T) {
	mockRepo := new(MockExportJobRepository)
	mockStorage := new(MockStorageClient) // Mock for S3/MinIO
	handler := NewExportHandler(mockRepo, nil) // Temporal not needed for download

	fileContent := []byte("csv,data,here")

	tests := []struct {
		name           string
		jobID          string
		setupMocks     func()
		expectedStatus int
		expectedHeader string
	}{
		{
			name:  "Success - Download File",
			jobID: "job-1",
			setupMocks: func() {
				job := &domain.ExportJobAggregate{
					ID:     "job-1",
					Status: "completed",
					FileKey: "files/job-1.csv",
				}
				mockRepo.On("Get", mock.Anything, "job-1").Return(job, nil)
				mockStorage.On("GetFile", "files/job-1.csv").Return(fileContent, "text/csv", nil)
			},
			expectedStatus: http.StatusOK,
			expectedHeader: "attachment; filename=job-1.csv",
		},
		{
			name:  "Failure - Job Not Found",
			jobID: "job-x",
			setupMocks: func() {
				mockRepo.On("Get", mock.Anything, "job-x").Return(nil, domain.ErrNotFound)
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:  "Failure - Job Pending",
			jobID: "job-2",
			setupMocks: func() {
				job := &domain.ExportJobAggregate{
					ID:     "job-2",
					Status: "processing",
				}
				mockRepo.On("Get", mock.Anything, "job-2").Return(job, nil)
			},
			expectedStatus: http.StatusAccepted, // Or 400 depending on spec
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()

			req := httptest.NewRequest("GET", "/exports/"+tc.jobID+"/download", nil)
			w := httptest.NewRecorder()

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tc.jobID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			handler.DownloadExport(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)
			if tc.expectedStatus == http.StatusOK {
				assert.Equal(t, tc.expectedHeader, w.Header().Get("Content-Disposition"))
			}
		})
	}
}
