package repository

import (
	"context"

	"github.com/sample-provider/buy-credit-api/internal/domain/entity"
)

type PartnerRepository interface {
	FindByClientID(ctx context.Context, clientID string) (*entity.Partner, error)
	FindByID(ctx context.Context, id string) (*entity.Partner, error)
}
