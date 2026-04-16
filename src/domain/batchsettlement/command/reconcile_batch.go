package command

// ReconcileBatchCmd is the command to finalize and reconcile a settlement batch.
// It verifies that the totals provided by the operator match the aggregated totals
// calculated by the system, ensuring data integrity before freezing the batch.
type ReconcileBatchCmd struct {
	BatchID             string
	ExpectedTotalDebits  int64
	ExpectedTotalCredits int64
}
