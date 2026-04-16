package event

// Event types
const (
	EventExportInitiated = "export.initiated"
)

// ExportInitiated represents the event when an export job is successfully initiated.
type ExportInitiated struct {
	JobID         string
	TargetDataset string
	FilterParams  map[string]interface{}
	Timestamp     int64
}
