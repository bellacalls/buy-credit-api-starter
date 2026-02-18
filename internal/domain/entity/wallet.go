package entity

type WalletStatus string

const (
	WalletStatusActive   WalletStatus = "ACTIVE"
	WalletStatusInactive WalletStatus = "INACTIVE"
	WalletStatusFrozen   WalletStatus = "FROZEN"
)

type Wallet struct {
	ID       string       `json:"id"`
	UserID   string       `json:"userId"`
	Currency string       `json:"currency"`
	Balance  string       `json:"balance"` // Using string for precision
	Status   WalletStatus `json:"status"`
}

func NewWallet(id, userID, currency, initialBalance string) *Wallet {
	return &Wallet{
		ID:       id,
		UserID:   userID,
		Currency: currency,
		Balance:  initialBalance,
		Status:   WalletStatusActive,
	}
}
