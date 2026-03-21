package domain

import (
	"time"

	"github.com/google/uuid"
)

type WebhookEventOriginal struct {
	ID          uuid.UUID
	WorkspaceID *uuid.UUID
	IntentID    *uuid.UUID
	Provider    string
	ExternalID  *string
	EventType   string
	Payload     []byte
	ReceivedAt  time.Time
}
