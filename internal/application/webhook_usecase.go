package application

import (
	"context"
	"fmt"

	"github.com/sample-provider/buy-credit-api/internal/domain/entity"
	"github.com/sample-provider/buy-credit-api/internal/domain/repository"
)

type WebhookUseCase struct {
	webhookRepo repository.WebhookRepository
}

func NewWebhookUseCase(webhookRepo repository.WebhookRepository) *WebhookUseCase {
	return &WebhookUseCase{
		webhookRepo: webhookRepo,
	}
}

type RegisterWebhookRequest struct {
	URL    string   `json:"url"`
	Events []string `json:"events"`
	Secret string   `json:"secret"`
}

type WebhookResponse struct {
	Webhook *entity.Webhook `json:"webhook"`
}

type WebhooksListResponse struct {
	Webhooks []*entity.Webhook `json:"webhooks"`
}

func (uc *WebhookUseCase) RegisterWebhook(ctx context.Context, partnerID string, req RegisterWebhookRequest) (*WebhookResponse, error) {
	// Validate request
	if req.URL == "" {
		return nil, fmt.Errorf("URL is required")
	}
	if len(req.Events) == 0 {
		return nil, fmt.Errorf("at least one event is required")
	}

	// Generate webhook ID
	webhookID := generateWebhookID()

	// Create webhook entity
	webhook := entity.NewWebhook(webhookID, partnerID, req.URL, req.Secret, req.Events)

	// Save to repository
	if err := uc.webhookRepo.Create(ctx, webhook); err != nil {
		return nil, fmt.Errorf("failed to create webhook: %w", err)
	}

	return &WebhookResponse{Webhook: webhook}, nil
}

func (uc *WebhookUseCase) GetWebhooks(ctx context.Context, partnerID string) (*WebhooksListResponse, error) {
	webhooks, err := uc.webhookRepo.FindByPartnerID(ctx, partnerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get webhooks: %w", err)
	}

	if webhooks == nil {
		webhooks = []*entity.Webhook{}
	}

	return &WebhooksListResponse{Webhooks: webhooks}, nil
}

func generateWebhookID() string {
	return fmt.Sprintf("whk_%d", generateRandomID())
}

func generateRandomID() int64 {
	return 1000000000 + int64(len("temp"))
}
