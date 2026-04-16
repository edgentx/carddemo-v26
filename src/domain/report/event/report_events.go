package event

const (
	// ReportRequestedEventType is the type for the event emitted when a report is requested.
	ReportRequestedEventType = "com.carddemo.report.requested"
)

// ReportRequested is emitted when a report generation is successfully requested.
type ReportRequested struct {
	ReportID   string                 `json:"report_id"`
	ConfigID   string                 `json:"config_id"`
	Format     string                 `json:"format"`
	Parameters map[string]interface{} `json:"parameters"`
}
