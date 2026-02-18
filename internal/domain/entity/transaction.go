package entity

import "time"

type TransactionStatus string
type TransactionType string

const (
	TransactionStatusPending TransactionStatus = "PENDING"
	TransactionStatusSuccess TransactionStatus = "SUCCESS"
	TransactionStatusFailed  TransactionStatus = "FAILED"

	TransactionTypeCreditPurchase TransactionType = "CREDIT_PURCHASE"
)

type Transaction struct {
	ID              string            `json:"id"`
	WalletID        string            `json:"walletId"`
	UserID          string            `json:"userId"`
	PartnerWalletID string            `json:"partnerWalletId"`
	Amount          string            `json:"amount"`
	Currency        string            `json:"currency"`
	Status          TransactionStatus `json:"status"`
	Type            TransactionType   `json:"type"`
	Metadata        map[string]string `json:"metadata,omitempty"`
	FailureReason   string            `json:"failureReason,omitempty"`
	CreatedAt       time.Time         `json:"createdAt"`
	UpdatedAt       time.Time         `json:"updatedAt"`
	CompletedAt     *time.Time        `json:"completedAt,omitempty"`
}

func NewTransaction(id, walletID, userID, partnerWalletID, amount, currency string, metadata map[string]string) *Transaction {
	now := time.Now()
	return &Transaction{
		ID:              id,
		WalletID:        walletID,
		UserID:          userID,
		PartnerWalletID: partnerWalletID,
		Amount:          amount,
		Currency:        currency,
		Status:          TransactionStatusPending,
		Type:            TransactionTypeCreditPurchase,
		Metadata:        metadata,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}

func (t *Transaction) MarkSuccess() {
	now := time.Now()
	t.Status = TransactionStatusSuccess
	t.UpdatedAt = now
	t.CompletedAt = &now
}

func (t *Transaction) MarkFailed(reason string) {
	now := time.Now()
	t.Status = TransactionStatusFailed
	t.FailureReason = reason
	t.UpdatedAt = now
	t.CompletedAt = &now
}
