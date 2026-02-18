package repository

import (
	"context"

	"github.com/sample-provider/buy-credit-api/internal/domain/entity"
)

type TransactionRepository interface {
	Create(ctx context.Context, transaction *entity.Transaction) error
	FindByID(ctx context.Context, id string) (*entity.Transaction, error)
	FindByIdempotencyKey(ctx context.Context, key string) (*entity.Transaction, error)
	Update(ctx context.Context, transaction *entity.Transaction) error
	StoreIdempotencyKey(ctx context.Context, key, transactionID string) error
}
