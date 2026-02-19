package handler

import (
	"encoding/json"
	"net/http"

	"github.com/sample-provider/buy-credit-api/internal/application"
	"github.com/sample-provider/buy-credit-api/internal/infrastructure/http/middleware"
	"github.com/sample-provider/buy-credit-api/internal/infrastructure/http/response"
)

type WebhookHandler struct {
	webhookUseCase *application.WebhookUseCase
}

func NewWebhookHandler(webhookUseCase *application.WebhookUseCase) *WebhookHandler {
	return &WebhookHandler{
		webhookUseCase: webhookUseCase,
	}
}

func (h *WebhookHandler) RegisterWebhook(w http.ResponseWriter, r *http.Request) {
	var req application.RegisterWebhookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	// Get partner ID from context (set by auth middleware)
	partnerID := middleware.GetPartnerID(r.Context())
	if partnerID == "" {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "Partner ID not found in context")
		return
	}

	webhookResp, err := h.webhookUseCase.RegisterWebhook(r.Context(), partnerID, req)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "WEBHOOK_CREATION_FAILED", err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, webhookResp)
}

func (h *WebhookHandler) GetWebhooks(w http.ResponseWriter, r *http.Request) {
	// Get partner ID from context (set by auth middleware)
	partnerID := middleware.GetPartnerID(r.Context())
	if partnerID == "" {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "Partner ID not found in context")
		return
	}

	webhooks, err := h.webhookUseCase.GetWebhooks(r.Context(), partnerID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "FAILED_TO_GET_WEBHOOKS", err.Error())
		return
	}

	response.JSON(w, http.StatusOK, webhooks)
}
