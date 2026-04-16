package command

// RequestReportCmd is a command to request a new report generation.
type RequestReportCmd struct {
	ReportID   string
	ConfigID   string
	Format     string
	Parameters map[string]interface{}
}
