package dto

import "time"

// ReportResponse represents the API response for a Report.
type ReportResponse struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}
