package command

// InitiateExportCmd is the command to trigger a bulk extraction of domain data.
type InitiateExportCmd struct {
	AggregateID     string
	TargetDataset  string
	FilterParams   map[string]interface{}
	UpstreamExists bool // Internal state/flag for validation (simulated)
}
