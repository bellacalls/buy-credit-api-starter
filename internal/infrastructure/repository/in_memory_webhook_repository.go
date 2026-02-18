package repository

import (
	"context"
	"errors"
	"sync"

	"github.com/sample-provider/buy-credit-api/internal/domain/entity"
	"github.com/sample-provider/buy-credit-api/internal/domain/repository"
)

type InMemoryWebhookRepository struct {
	mu       sync.RWMutex
	webhooks map[string]*entity.Webhook
}

func NewInMemoryWebhookRepository() repository.WebhookRepository {
	return &InMemoryWebhookRepository{
		webhooks: make(map[string]*entity.Webhook),
	}
}

func (r *InMemoryWebhookRepository) Create(ctx context.Context, webhook *entity.Webhook) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.webhooks[webhook.ID]; exists {
		return errors.New("webhook already exists")
	}

	r.webhooks[webhook.ID] = webhook
	return nil
}

func (r *InMemoryWebhookRepository) FindByPartnerID(ctx context.Context, partnerID string) ([]*entity.Webhook, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var webhooks []*entity.Webhook
	for _, webhook := range r.webhooks {
		if webhook.PartnerID == partnerID {
			webhooks = append(webhooks, webhook)
		}
	}

	return webhooks, nil
}

func (r *InMemoryWebhookRepository) FindByID(ctx context.Context, id string) (*entity.Webhook, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	webhook, exists := r.webhooks[id]
	if !exists {
		return nil, errors.New("webhook not found")
	}

	return webhook, nil
}
