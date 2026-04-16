package event

// Event types
const (
	EventExportInitiated = "export.initiated"
	EventExportCompleted = "export.completed"
)

// ExportInitiated represents the event when an export job is successfully initiated.
type ExportInitiated struct {
	JobID         string
	TargetDataset string
	FilterParams  map[string]interface{}
	Timestamp     int64
}

// ExportCompleted represents the event when an export job is successfully finished.
type ExportCompleted struct {
	JobID        string
	RecordCount  int64
	ManifestData string // The generated manifest content
	Timestamp    int64
}
