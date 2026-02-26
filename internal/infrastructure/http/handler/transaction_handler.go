package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sample-provider/buy-credit-api/internal/application"
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
	if req.UserID == "" || req.Amount == 0 || req.Currency == "" {
		response.Error(w, http.StatusBadRequest, "MISSING_FIELDS", "userId, amount, and currency are required")
		return
	}

	txnResp, err := h.transactionUseCase.CreateTransaction(r.Context(), req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		code := "INTERNAL_ERROR"

		switch err.Error() {
		case "invalid amount":
			statusCode = http.StatusBadRequest
			code = "INVALID_AMOUNT"
		}

		response.Error(w, statusCode, code, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, txnResp)
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
