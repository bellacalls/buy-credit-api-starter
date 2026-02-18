package repository

import (
	"context"

	"github.com/sample-provider/buy-credit-api/internal/domain/entity"
)

type WebhookRepository interface {
	Create(ctx context.Context, webhook *entity.Webhook) error
	FindByPartnerID(ctx context.Context, partnerID string) ([]*entity.Webhook, error)
	FindByID(ctx context.Context, id string) (*entity.Webhook, error)
}
