package command

// CompleteExportCmd is the command to mark a bulk export process as successfully finished.
// It validates the generated manifest and finalizes the job.
type CompleteExportCmd struct {
	AggregateID  string
	ExportID     string // Valid export_id
	RecordCount  int64   // Valid record count
	ManifestData string // Valid manifest (e.g., JSON string or File Handle)
}
