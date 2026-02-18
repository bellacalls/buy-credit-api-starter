package repository

import (
	"context"
	"errors"
	"sync"

	"github.com/sample-provider/buy-credit-api/internal/domain/entity"
	"github.com/sample-provider/buy-credit-api/internal/domain/repository"
)

type InMemoryTransactionRepository struct {
	mu              sync.RWMutex
	transactions    map[string]*entity.Transaction
	idempotencyKeys map[string]string // idempotencyKey -> transactionID
}

func NewInMemoryTransactionRepository() repository.TransactionRepository {
	return &InMemoryTransactionRepository{
		transactions:    make(map[string]*entity.Transaction),
		idempotencyKeys: make(map[string]string),
	}
}

func (r *InMemoryTransactionRepository) Create(ctx context.Context, transaction *entity.Transaction) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.transactions[transaction.ID]; exists {
		return errors.New("transaction already exists")
	}

	r.transactions[transaction.ID] = transaction
	return nil
}

func (r *InMemoryTransactionRepository) FindByID(ctx context.Context, id string) (*entity.Transaction, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	transaction, exists := r.transactions[id]
	if !exists {
		return nil, errors.New("transaction not found")
	}

	return transaction, nil
}

func (r *InMemoryTransactionRepository) FindByIdempotencyKey(ctx context.Context, key string) (*entity.Transaction, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	transactionID, exists := r.idempotencyKeys[key]
	if !exists {
		return nil, errors.New("transaction not found")
	}

	return r.transactions[transactionID], nil
}

func (r *InMemoryTransactionRepository) Update(ctx context.Context, transaction *entity.Transaction) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.transactions[transaction.ID]; !exists {
		return errors.New("transaction not found")
	}

	r.transactions[transaction.ID] = transaction
	return nil
}

func (r *InMemoryTransactionRepository) StoreIdempotencyKey(ctx context.Context, key, transactionID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.idempotencyKeys[key] = transactionID
	return nil
}
