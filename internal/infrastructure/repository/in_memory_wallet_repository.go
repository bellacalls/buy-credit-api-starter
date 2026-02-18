package repository

import (
	"context"
	"errors"
	"sync"

	"github.com/sample-provider/buy-credit-api/internal/domain/entity"
	"github.com/sample-provider/buy-credit-api/internal/domain/repository"
)

type InMemoryWalletRepository struct {
	mu      sync.RWMutex
	wallets map[string]*entity.Wallet
}

func NewInMemoryWalletRepository() repository.WalletRepository {
	repo := &InMemoryWalletRepository{
		wallets: make(map[string]*entity.Wallet),
	}

	// Seed with sample data
	repo.seedData()
	return repo
}

func (r *InMemoryWalletRepository) seedData() {
	// Customer wallet
	wallet1 := entity.NewWallet("wlt_usd_abc123", "usr_123", "USD", "1500.50")
	r.wallets[wallet1.ID] = wallet1

	// Partner wallet (Bella Mobile)
	wallet2 := entity.NewWallet("wlt_partner_bella", "partner_bella", "USD", "0.00")
	r.wallets[wallet2.ID] = wallet2
}

func (r *InMemoryWalletRepository) FindByID(ctx context.Context, id string) (*entity.Wallet, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	wallet, exists := r.wallets[id]
	if !exists {
		return nil, errors.New("wallet not found")
	}

	return wallet, nil
}

func (r *InMemoryWalletRepository) FindByUserID(ctx context.Context, userID string) ([]*entity.Wallet, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var wallets []*entity.Wallet
	for _, wallet := range r.wallets {
		if wallet.UserID == userID {
			wallets = append(wallets, wallet)
		}
	}

	return wallets, nil
}

func (r *InMemoryWalletRepository) Update(ctx context.Context, wallet *entity.Wallet) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.wallets[wallet.ID]; !exists {
		return errors.New("wallet not found")
	}

	r.wallets[wallet.ID] = wallet
	return nil
}
