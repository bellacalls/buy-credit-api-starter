package entity

import "time"

type WebhookStatus string

const (
	WebhookStatusActive   WebhookStatus = "ACTIVE"
	WebhookStatusInactive WebhookStatus = "INACTIVE"
)

type Webhook struct {
	ID        string        `json:"id"`
	PartnerID string        `json:"partnerId"`
	URL       string        `json:"url"`
	Events    []string      `json:"events"`
	Secret    string        `json:"secret"`
	Status    WebhookStatus `json:"status"`
	CreatedAt time.Time     `json:"createdAt"`
	UpdatedAt time.Time     `json:"updatedAt"`
}

func NewWebhook(id, partnerID, url, secret string, events []string) *Webhook {
	now := time.Now()
	return &Webhook{
		ID:        id,
		PartnerID: partnerID,
		URL:       url,
		Events:    events,
		Secret:    secret,
		Status:    WebhookStatusActive,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

type WebhookEvent struct {
	EventID   string                 `json:"eventId"`
	EventType string                 `json:"eventType"`
	Timestamp time.Time              `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
}

func NewWebhookEvent(eventType string, data map[string]interface{}) *WebhookEvent {
	return &WebhookEvent{
		EventID:   generateEventID(),
		EventType: eventType,
		Timestamp: time.Now(),
		Data:      data,
	}
}

func generateEventID() string {
	return "evt_" + time.Now().Format("20060102150405")
}
