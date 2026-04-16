package dto

import "time"

// ExportJobResponse represents the API response for an Export Job.
type ExportJobResponse struct {
	ID        string    `json:"id"`
	ReportID  string    `json:"report_id"`
	Status    string    `json:"status"`
	FileKey   string    `json:"file_key,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
