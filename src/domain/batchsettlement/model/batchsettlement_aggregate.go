package model

import (
	"errors"

	"github.com/carddemo/project/src/domain/batchsettlement/command"
	"github.com/carddemo/project/src/domain/batchsettlement/event"
	"github.com/carddemo/project/src/domain/shared"
)

var (
	// ErrUncommittedTransactions is returned when opening a batch with pending/uncommitted transactions.
	ErrUncommittedTransactions = errors.New("batch.settlement.uncommitted_transactions")

	// ErrBalanceMismatch is returned when opening a batch that is not financially balanced.
	ErrBalanceMismatch = errors.New("batch.settlement.balance_mismatch")
)

// BatchSettlement represents the BatchSettlement aggregate.
type BatchSettlement struct {
	shared.AggregateRoot
	ID                 string
	HasUncommittedTxns bool // Invariant: Flag indicating presence of pending transactions
	IsBalanced         bool // Invariant: Flag indicating financial balance (Debits == Credits)
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

// GetID returns the aggregate ID.
func (b *BatchSettlement) GetID() string {
	return b.ID
}

// ID satisfies the shared.Aggregate interface.
func (b *BatchSettlement) ID() string {
	return b.ID
}
