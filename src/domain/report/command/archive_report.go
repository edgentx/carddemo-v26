package command

// ArchiveReportCmd is a command to persist a generated report artifact to long-term object storage.
type ArchiveReportCmd struct {
	ReportID         string
	StorageLocation  string
	// ArtifactChecksum could be included for integrity verification
}
