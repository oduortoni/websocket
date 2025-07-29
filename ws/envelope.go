package ws

import (
	"time"
)

type Envelope struct {
	ID        Identity               `json:"id"`
	ClientID  Identity               `json:"client_id"`
	Type      string                 `json:"type"`
	Payload   map[string]interface{} `json:"payload"`
	Timestamp time.Time              `json:"timestamp"`
	Delivered *time.Time             `json:"delivered"`
}

type EnvelopePersister interface {
	SaveEnvelope(e Envelope) error
	ConfirmDelivery(envelopeID Identity, clientID Identity) error
}
