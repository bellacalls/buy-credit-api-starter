package application

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/sample-provider/buy-credit-api/internal/domain/entity"
	"github.com/sample-provider/buy-credit-api/internal/domain/repository"
)

type TransactionUseCase struct {
	transactionRepo repository.TransactionRepository
}

type CreateTransactionRequest struct {
	UserID   string  `json:"userId"`
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

type TransactionResponse struct {
	ID        string                   `json:"transactionId"`
	UserID    string                   `json:"userId"`
	Currency  string                   `json:"currency"`
	Amount    float64                  `json:"amount"`
	Status    entity.TransactionStatus `json:"status"`
	Timestamp string                   `json:"timestamp"`
}

func NewTransactionUseCase(
	transactionRepo repository.TransactionRepository,
) *TransactionUseCase {
	return &TransactionUseCase{
		transactionRepo: transactionRepo,
	}
}

func (uc *TransactionUseCase) CreateTransaction(ctx context.Context, req CreateTransactionRequest) (*TransactionResponse, error) {
	// Validate amount
	if req.Amount <= 0 {
		return nil, errors.New("invalid amount")
	}

	// Create transaction
	txnID := fmt.Sprintf("txn_%s", uuid.New().String()[:8])
	transaction := entity.NewTransaction(
		txnID,
		req.UserID,
		req.Currency,
		req.Amount,
	)

	if err := uc.transactionRepo.Create(ctx, transaction); err != nil {
		return nil, err
	}

	return &TransactionResponse{
		ID:        transaction.ID,
		UserID:    transaction.UserID,
		Currency:  transaction.Currency,
		Amount:    transaction.Amount,
		Status:    transaction.Status,
		Timestamp: transaction.Timestamp.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

func (uc *TransactionUseCase) GetTransaction(ctx context.Context, transactionID string) (*TransactionResponse, error) {
	transaction, err := uc.transactionRepo.FindByID(ctx, transactionID)
	if err != nil {
		return nil, errors.New("transaction not found")
	}

	return &TransactionResponse{
		ID:        transaction.ID,
		UserID:    transaction.UserID,
		Currency:  transaction.Currency,
		Amount:    transaction.Amount,
		Status:    transaction.Status,
		Timestamp: transaction.Timestamp.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

func (uc *TransactionUseCase) UpdateTransactionStatus(ctx context.Context, transactionID string, status entity.TransactionStatus) error {
	transaction, err := uc.transactionRepo.FindByID(ctx, transactionID)
	if err != nil {
		return errors.New("transaction not found")
	}

	switch status {
	case entity.TransactionStatusSuccessful:
		transaction.MarkSuccessful()
	case entity.TransactionStatusFailed:
		transaction.MarkFailed()
	default:
		return errors.New("invalid status")
	}

	return uc.transactionRepo.Update(ctx, transaction)
}
