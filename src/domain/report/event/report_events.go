package event

const (
	// ReportRequestedEventType is the type for the event emitted when a report is requested.
	ReportRequestedEventType = "com.carddemo.report.requested"
	// ReportArchivedEventType is the type for the event emitted when a report is successfully archived.
	ReportArchivedEventType = "com.carddemo.report.archived"
)

// ReportRequested is emitted when a report generation is successfully requested.
type ReportRequested struct {
	ReportID   string                 `json:"report_id"`
	ConfigID   string                 `json:"config_id"`
	Format     string                 `json:"format"`
	Parameters map[string]interface{} `json:"parameters"`
}

// ReportArchived is emitted when a report artifact is successfully persisted to storage.
type ReportArchived struct {
	ReportID        string `json:"report_id"`
	StorageLocation string `json:"storage_location"`
}
