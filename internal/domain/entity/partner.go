package entity

import "time"

type PartnerStatus string

const (
	PartnerStatusActive   PartnerStatus = "ACTIVE"
	PartnerStatusInactive PartnerStatus = "INACTIVE"
)

type Partner struct {
	ID           string        `json:"id"`
	Name         string        `json:"name"`
	ClientID     string        `json:"clientId"`
	ClientSecret string        `json:"-"` // Never expose in JSON
	WalletID     string        `json:"walletId"`
	Status       PartnerStatus `json:"status"`
	CreatedAt    time.Time     `json:"createdAt"`
	UpdatedAt    time.Time     `json:"updatedAt"`
}

func NewPartner(id, name, clientID, clientSecret, walletID string) *Partner {
	now := time.Now()
	return &Partner{
		ID:           id,
		Name:         name,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		WalletID:     walletID,
		Status:       PartnerStatusActive,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}
