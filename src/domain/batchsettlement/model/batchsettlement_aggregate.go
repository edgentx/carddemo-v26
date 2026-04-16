package model

import (
	"errors"

	"github.com/carddemo/project/src/domain/batchsettlement/command"
	"github.com/carddemo/project/src/domain/batchsettlement/event"
	"github.com/carddemo/project/src/domain/shared"
)

var (
	// ErrUncommittedTransactions is returned when reconciling/finalizing a batch with pending transactions.
	ErrUncommittedTransactions = errors.New("batch.settlement.uncommitted_transactions")

	// ErrBalanceMismatch is returned when opening a batch that is not financially balanced.
	ErrBalanceMismatch = errors.New("batch.settlement.balance_mismatch")

	// ErrInvalidReconciliationTotals is returned when expected totals do not match aggregate state.
	ErrInvalidReconciliationTotals = errors.New("batch.settlement.invalid_reconciliation_totals")
)

// BatchSettlement represents the BatchSettlement aggregate.
type BatchSettlement struct {
	shared.AggregateRoot
	ID                 string
	HasUncommittedTxns bool // Invariant: Flag indicating presence of pending transactions
	IsBalanced         bool // Invariant: Flag indicating financial balance (Debits == Credits)
	TotalDebits        int64 // Current total of debit transactions in the batch
	TotalCredits       int64 // Current total of credit transactions in the batch
}

// NewBatchSettlement creates a new BatchSettlement instance.
func NewBatchSettlement(id string) *BatchSettlement {
	return &BatchSettlement{
		ID:         id,
		IsBalanced: true, // Default state for a clean slate
	}
}

// Execute handles commands for the BatchSettlement aggregate.
func (b *BatchSettlement) Execute(cmd interface{}) ([]shared.DomainEvent, error) {
	switch c := cmd.(type) {
	case command.OpenBatchCmd:
		return b.handleOpenBatch(c)
	case command.ReconcileBatchCmd:
		return b.handleReconcileBatch(c)
	default:
		return nil, shared.ErrUnknownCommand
	}
}

// handleOpenBatch processes the OpenBatchCmd command.
func (b *BatchSettlement) handleOpenBatch(cmd command.OpenBatchCmd) ([]shared.DomainEvent, error) {
	// 1. Check Invariants
	if b.HasUncommittedTxns {
		return nil, ErrUncommittedTransactions
	}

	if !b.IsBalanced {
		return nil, ErrBalanceMismatch
	}

	// 2. Generate Event
	e := event.BatchOpened{
		DomainEventMeta:   shared.NewDomainEventMeta(b.ID),
		SettlementDate:    cmd.SettlementDate,
		OperationalRegion: cmd.OperationalRegion,
	}

	return []shared.DomainEvent{e}, nil
}

// handleReconcileBatch processes the ReconcileBatchCmd command.
func (b *BatchSettlement) handleReconcileBatch(cmd command.ReconcileBatchCmd) ([]shared.DomainEvent, error) {
	// 1. Validate Invariants: Cannot reconcile if transactions are pending
	if b.HasUncommittedTxns {
		return nil, ErrUncommittedTransactions
	}

	// 2. Validate Totals: Strict comparison against current state.
	// We do not perform internal balancing logic (Debits == Credits) here;
	// we enforce that the caller's expectation matches the aggregate's calculated truth.
	if b.TotalDebits != cmd.ExpectedTotalDebits || b.TotalCredits != cmd.ExpectedTotalCredits {
		return nil, ErrInvalidReconciliationTotals
	}

	// 3. Generate Event
	e := event.BatchReconciled{
		DomainEventMeta: shared.NewDomainEventMeta(b.ID),
		TotalDebits:     b.TotalDebits,
		TotalCredits:    b.TotalCredits,
	}

	return []shared.DomainEvent{e}, nil
}

// GetID returns the aggregate ID.
func (b *BatchSettlement) GetID() string {
	return b.ID
}

// ID satisfies the shared.Aggregate interface.
func (b *BatchSettlement) ID() string {
	return b.ID
}
