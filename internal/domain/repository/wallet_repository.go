package repository

import (
	"context"

	"github.com/sample-provider/buy-credit-api/internal/domain/entity"
)

type WalletRepository interface {
	FindByID(ctx context.Context, id string) (*entity.Wallet, error)
	FindByUserID(ctx context.Context, userID string) ([]*entity.Wallet, error)
	Update(ctx context.Context, wallet *entity.Wallet) error
}
