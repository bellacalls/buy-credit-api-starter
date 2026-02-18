package application

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/sample-provider/buy-credit-api/internal/domain/entity"
	"github.com/sample-provider/buy-credit-api/internal/domain/repository"
)

type TransactionUseCase struct {
	transactionRepo repository.TransactionRepository
	walletRepo      repository.WalletRepository
	webhookRepo     repository.WebhookRepository
}

type CreateTransactionRequest struct {
	UserID          string            `json:"userId"`
	WalletID        string            `json:"walletId"`
	Amount          string            `json:"amount"`
	Currency        string            `json:"currency"`
	Metadata        map[string]string `json:"metadata"`
	IdempotencyKey  string            `json:"idempotencyKey,omitempty"`
	PartnerID       string            `json:"-"`
	PartnerWalletID string            `json:"-"`
}

type TransactionResponse struct {
	Transaction *entity.Transaction `json:"transaction"`
}

func NewTransactionUseCase(
	transactionRepo repository.TransactionRepository,
	walletRepo repository.WalletRepository,
	webhookRepo repository.WebhookRepository,
) *TransactionUseCase {
	return &TransactionUseCase{
		transactionRepo: transactionRepo,
		walletRepo:      walletRepo,
		webhookRepo:     webhookRepo,
	}
}

func (uc *TransactionUseCase) CreateTransaction(ctx context.Context, req CreateTransactionRequest) (*TransactionResponse, error) {
	// Check idempotency
	if req.IdempotencyKey != "" {
		existingTxn, err := uc.transactionRepo.FindByIdempotencyKey(ctx, req.IdempotencyKey)
		if err == nil && existingTxn != nil {
			return &TransactionResponse{Transaction: existingTxn}, nil
		}
	}

	// Validate wallet
	wallet, err := uc.walletRepo.FindByID(ctx, req.WalletID)
	if err != nil {
		return nil, errors.New("wallet not found")
	}

	if wallet.UserID != req.UserID {
		return nil, errors.New("wallet does not belong to user")
	}

	if wallet.Status != entity.WalletStatusActive {
		return nil, errors.New("wallet is not active")
	}

	// Check balance
	balance, _ := strconv.ParseFloat(wallet.Balance, 64)
	amount, _ := strconv.ParseFloat(req.Amount, 64)

	if balance < amount {
		return nil, errors.New("insufficient balance")
	}

	// Create transaction
	txnID := fmt.Sprintf("txn_purchase_%s", uuid.New().String()[:8])
	transaction := entity.NewTransaction(
		txnID,
		req.WalletID,
		req.UserID,
		req.PartnerWalletID,
		req.Amount,
		req.Currency,
		req.Metadata,
	)

	if err := uc.transactionRepo.Create(ctx, transaction); err != nil {
		return nil, err
	}

	// Store idempotency key
	if req.IdempotencyKey != "" {
		uc.transactionRepo.StoreIdempotencyKey(ctx, req.IdempotencyKey, transaction.ID)
	}

	// Process transaction asynchronously (simulated here)
	go uc.processTransaction(ctx, transaction, wallet, amount)

	return &TransactionResponse{Transaction: transaction}, nil
}

func (uc *TransactionUseCase) GetTransaction(ctx context.Context, transactionID string) (*TransactionResponse, error) {
	transaction, err := uc.transactionRepo.FindByID(ctx, transactionID)
	if err != nil {
		return nil, errors.New("transaction not found")
	}

	return &TransactionResponse{Transaction: transaction}, nil
}

func (uc *TransactionUseCase) processTransaction(ctx context.Context, txn *entity.Transaction, wallet *entity.Wallet, amount float64) {
	// Simulate processing
	// In production: deduct from customer wallet, credit partner wallet, call partner API, etc.

	// Update wallet balance
	currentBalance, _ := strconv.ParseFloat(wallet.Balance, 64)
	newBalance := currentBalance - amount
	wallet.Balance = fmt.Sprintf("%.2f", newBalance)
	uc.walletRepo.Update(ctx, wallet)

	// Mark transaction as success
	txn.MarkSuccess()
	uc.transactionRepo.Update(ctx, txn)

	// Send webhook notification (if webhooks are registered)
	// This would be done via a proper webhook delivery service in production
}
