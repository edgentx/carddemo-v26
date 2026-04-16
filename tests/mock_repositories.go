package tests

import (
	"errors"
	"sync"
	"time"

	"github.com/carddemo/project/src/domain/batchsettlement/model"
	"github.com/carddemo/project/src/domain/batchsettlement/repository"
	"github.com/carddemo/project/src/domain/transaction/model"
	"github.com/carddemo/project/src/domain/transaction/repository"
)

// MockTransactionRepositoryWithQuerySupport extends the basic mock to include query methods.
type MockTransactionRepositoryWithQuerySupport struct {
	mu            sync.RWMutex
	transactions  map[string]*model.Transaction
	indexesCreated bool
}

func NewMockTransactionRepositoryWithQuerySupport() *MockTransactionRepositoryWithQuerySupport {
	return &MockTransactionRepositoryWithQuerySupport{
		transactions: make(map[string]*model.Transaction),
	}
}

var _ repository.TransactionRepository = (*MockTransactionRepositoryWithQuerySupport)(nil)

func (m *MockTransactionRepositoryWithQuerySupport) Get(id string) (*model.Transaction, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.transactions[id], nil
}

func (m *MockTransactionRepositoryWithQuerySupport) Save(aggregate *model.Transaction) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.transactions[aggregate.ID] = aggregate
	return nil
}

func (m *MockTransactionRepositoryWithQuerySupport) List() ([]*model.Transaction, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var arr []*model.Transaction
	for _, v := range m.transactions {
		arr = append(arr, v)
	}
	return arr, nil
}

func (m *MockTransactionRepositoryWithQuerySupport) CreateIndexes() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.indexesCreated = true
	return nil
}

func (m *MockTransactionRepositoryWithQuerySupport) FindByCardAndDateRange(cardId string, start, end time.Time) ([]*model.Transaction, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var res []*model.Transaction
	for _, t := range m.transactions {
		if t.CardID == cardId && t.Timestamp.After(start) && t.Timestamp.Before(end) {
			res = append(res, t)
		}
	}
	return res, nil
}

func (m *MockTransactionRepositoryWithQuerySupport) FindByStatus(status string) ([]*model.Transaction, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var res []*model.Transaction
	for _, t := range m.transactions {
		if t.Status == status {
			res = append(res, t)
		}
	}
	return res, nil
}

// MockBatchSettlementRepositoryWithBulkSupport supports bulk operations and aggregation queries.
type MockBatchSettlementRepositoryWithBulkSupport struct {
	mu      sync.RWMutex
	batches map[string]*model.BatchSettlement
}

func NewMockBatchSettlementRepositoryWithBulkSupport() *MockBatchSettlementRepositoryWithBulkSupport {
	return &MockBatchSettlementRepositoryWithBulkSupport{
		batches: make(map[string]*model.BatchSettlement),
	}
}

var _ repository.BatchSettlementRepository = (*MockBatchSettlementRepositoryWithBulkSupport)(nil)

func (m *MockBatchSettlementRepositoryWithBulkSupport) Get(id string) (*model.BatchSettlement, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.batches[id], nil
}

func (m *MockBatchSettlementRepositoryWithBulkSupport) Save(aggregate *model.BatchSettlement) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.batches[aggregate.ID] = aggregate
	return nil
}

func (m *MockBatchSettlementRepositoryWithBulkSupport) List() ([]*model.BatchSettlement, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var arr []*model.BatchSettlement
	for _, v := range m.batches {
		arr = append(arr, v)
	}
	return arr, nil
}

func (m *MockBatchSettlementRepositoryWithBulkSupport) BulkInsertTransactions(transactions []*model.Transaction) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	// In a real scenario, this would insert to a Transaction collection.
	// For the mock, we can accept the call to verify interface compliance.
	if len(transactions) == 0 {
		return errors.New("cannot insert empty batch")
	}
	return nil
}

func (m *MockBatchSettlementRepositoryWithBulkSupport) GetSettlementAggregation(startDate, endDate time.Time) ([]*model.SettlementGroup, error) {
	// Return dummy data to simulate the aggregation pipeline structure
	return []*model.SettlementGroup{
		{MerchantID: "m1", Date: startDate, Count: 1, Total: 100.0},
	}, nil
}
