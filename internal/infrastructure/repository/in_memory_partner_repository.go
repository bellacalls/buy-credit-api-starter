package repository

import (
	"context"
	"errors"
	"sync"

	"github.com/sample-provider/buy-credit-api/internal/domain/entity"
	"github.com/sample-provider/buy-credit-api/internal/domain/repository"
)

type InMemoryPartnerRepository struct {
	mu       sync.RWMutex
	partners map[string]*entity.Partner
}

func NewInMemoryPartnerRepository() repository.PartnerRepository {
	repo := &InMemoryPartnerRepository{
		partners: make(map[string]*entity.Partner),
	}

	// Seed with sample partner data
	repo.seedData()
	return repo
}

func (r *InMemoryPartnerRepository) seedData() {
	partner := entity.NewPartner(
		"partner_bella",
		"Bella Mobile",
		"bella_mobile_prod",
		"secret_bella_123",
		"wlt_partner_bella",
	)
	r.partners[partner.ID] = partner
}

func (r *InMemoryPartnerRepository) FindByClientID(ctx context.Context, clientID string) (*entity.Partner, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, partner := range r.partners {
		if partner.ClientID == clientID {
			return partner, nil
		}
	}

	return nil, errors.New("partner not found")
}

func (r *InMemoryPartnerRepository) FindByID(ctx context.Context, id string) (*entity.Partner, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	partner, exists := r.partners[id]
	if !exists {
		return nil, errors.New("partner not found")
	}

	return partner, nil
}
