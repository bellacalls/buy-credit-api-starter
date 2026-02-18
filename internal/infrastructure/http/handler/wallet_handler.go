package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sample-provider/buy-credit-api/internal/application"
	"github.com/sample-provider/buy-credit-api/internal/infrastructure/http/response"
)

type WalletHandler struct {
	walletUseCase *application.WalletUseCase
}

func NewWalletHandler(walletUseCase *application.WalletUseCase) *WalletHandler {
	return &WalletHandler{
		walletUseCase: walletUseCase,
	}
}

func (h *WalletHandler) GetWallets(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	if userID == "" {
		response.Error(w, http.StatusBadRequest, "MISSING_USER_ID", "X-User-ID header is required")
		return
	}

	wallets, err := h.walletUseCase.GetUserWallets(r.Context(), userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch wallets")
		return
	}

	response.JSON(w, http.StatusOK, wallets)
}
