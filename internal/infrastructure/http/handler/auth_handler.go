package handler

import (
	"encoding/json"
	"net/http"

	"github.com/sample-provider/buy-credit-api/internal/application"
	"github.com/sample-provider/buy-credit-api/internal/infrastructure/http/response"
)

type AuthHandler struct {
	authUseCase *application.AuthUseCase
}

func NewAuthHandler(authUseCase *application.AuthUseCase) *AuthHandler {
	return &AuthHandler{
		authUseCase: authUseCase,
	}
}

func (h *AuthHandler) CreateToken(w http.ResponseWriter, r *http.Request) {
	var req application.AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	if req.APIKey == "" || req.APISecret == "" {
		response.Error(w, http.StatusBadRequest, "MISSING_FIELDS", "apiKey and apiSecret are required")
		return
	}

	authResp, err := h.authUseCase.Authenticate(r.Context(), req)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "INVALID_CREDENTIALS", "Invalid client credentials")
		return
	}

	response.JSON(w, http.StatusOK, authResp)
}
