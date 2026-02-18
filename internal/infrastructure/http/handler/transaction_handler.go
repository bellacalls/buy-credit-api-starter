package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sample-provider/buy-credit-api/internal/application"
	"github.com/sample-provider/buy-credit-api/internal/infrastructure/http/middleware"
	"github.com/sample-provider/buy-credit-api/internal/infrastructure/http/response"
)

type TransactionHandler struct {
	transactionUseCase *application.TransactionUseCase
}

func NewTransactionHandler(transactionUseCase *application.TransactionUseCase) *TransactionHandler {
	return &TransactionHandler{
		transactionUseCase: transactionUseCase,
	}
}

func (h *TransactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var req application.CreateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	// Validate required fields
	if req.UserID == "" || req.WalletID == "" || req.Amount == "" || req.Currency == "" {
		response.Error(w, http.StatusBadRequest, "MISSING_FIELDS", "userId, walletId, amount, and currency are required")
		return
	}

	partnerID := middleware.GetPartnerID(r.Context())
	if partnerID == "" {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid authentication")
		return
	}

	// Get partner wallet ID (in production, fetch from partner entity)
	req.PartnerID = partnerID
	req.PartnerWalletID = "wlt_partner_bella" // In production, get from partner entity
	req.IdempotencyKey = r.Header.Get("Idempotency-Key")

	txnResp, err := h.transactionUseCase.CreateTransaction(r.Context(), req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		code := "INTERNAL_ERROR"

		switch err.Error() {
		case "wallet not found":
			statusCode = http.StatusNotFound
			code = "WALLET_NOT_FOUND"
		case "wallet does not belong to user":
			statusCode = http.StatusForbidden
			code = "FORBIDDEN"
		case "wallet is not active":
			statusCode = http.StatusBadRequest
			code = "WALLET_INACTIVE"
		case "insufficient balance":
			statusCode = http.StatusBadRequest
			code = "INSUFFICIENT_BALANCE"
		}

		response.Error(w, statusCode, code, err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, txnResp)
}

func (h *TransactionHandler) GetTransaction(w http.ResponseWriter, r *http.Request) {
	transactionID := chi.URLParam(r, "transactionId")
	if transactionID == "" {
		response.Error(w, http.StatusBadRequest, "MISSING_TRANSACTION_ID", "Transaction ID is required")
		return
	}

	txnResp, err := h.transactionUseCase.GetTransaction(r.Context(), transactionID)
	if err != nil {
		response.Error(w, http.StatusNotFound, "TRANSACTION_NOT_FOUND", "Transaction not found")
		return
	}

	response.JSON(w, http.StatusOK, txnResp)
}
