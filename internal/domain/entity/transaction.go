package entity

import "time"

type TransactionStatus string
type TransactionType string

const (
	TransactionStatusPending    TransactionStatus = "PENDING"
	TransactionStatusSuccessful TransactionStatus = "SUCCESSFUL"
	TransactionStatusFailed     TransactionStatus = "FAILED"

	TransactionTypeCreditPurchase TransactionType = "CREDIT_PURCHASE"
)

type Transaction struct {
	ID        string            `json:"transactionId"`
	UserID    string            `json:"userId"`
	Amount    float64           `json:"amount"`
	Currency  string            `json:"currency"`
	Status    TransactionStatus `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
}

func NewTransaction(id, userID, currency string, amount float64) *Transaction {
	return &Transaction{
		ID:        id,
		UserID:    userID,
		Amount:    amount,
		Currency:  currency,
		Status:    TransactionStatusPending,
		Timestamp: time.Now(),
	}
}

func (t *Transaction) MarkSuccessful() {
	t.Status = TransactionStatusSuccessful
	t.Timestamp = time.Now()
}

func (t *Transaction) MarkFailed() {
	t.Status = TransactionStatusFailed
	t.Timestamp = time.Now()
}
