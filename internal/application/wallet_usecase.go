package application

import (
	"context"

	"github.com/sample-provider/buy-credit-api/internal/domain/entity"
	"github.com/sample-provider/buy-credit-api/internal/domain/repository"
)

type WalletUseCase struct {
	walletRepo repository.WalletRepository
}

type WalletResponse struct {
	Wallets []*entity.Wallet `json:"wallets"`
}

func NewWalletUseCase(walletRepo repository.WalletRepository) *WalletUseCase {
	return &WalletUseCase{
		walletRepo: walletRepo,
	}
}

func (uc *WalletUseCase) GetUserWallets(ctx context.Context, userID string) (*WalletResponse, error) {
	wallets, err := uc.walletRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &WalletResponse{
		Wallets: wallets,
	}, nil
}
