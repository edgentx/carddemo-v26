package service

// Commands for BatchSettlement Aggregate

type OpenBatchCmd struct {
	BatchID string
	Name    string
}

type ReconcileBatchCmd struct {
	BatchID string
}
