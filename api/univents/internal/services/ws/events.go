package ws

import (
	"encoding/json"
	"time"

	"univents/models"

	"github.com/google/uuid"
)

// Frame types (D9): the snapshot restores context on connect, then live
// frames until a terminal frame (confirmed/expired/cancelled) closes the
// socket. Every frame is `{"type": "...", "payload": {...}}`.
const (
	frameSnapshot  = "purchase.snapshot"
	frameIntent    = "intent.updated"
	frameConfirmed = "purchase.confirmed"
	frameExpired   = "purchase.expired"
	frameCancelled = "purchase.cancelled"
)

// intent statuses the hub emits in intent.updated. Mirrors the payssage
// vocabulary for the states univents ever sees (helpers.go NormalizeStatus);
// the hub derives them from the purchase status the notification carries.
const (
	intentSucceeded = "succeeded"
	intentCancelled = "cancelled"
)

// frame is the wire envelope. Payload is pre-marshaled so the frame builder
// never fails on payload shape.
type frame struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

func marshalFrame(typ string, payload any) ([]byte, error) {
	raw, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return json.Marshal(frame{Type: typ, Payload: raw})
}

// snapshotPayload mirrors the split-5 resume (getCheckout) shape so the
// front treats resume == snapshot == checkout uniformly. `intent_status` is
// best-effort: present only when the purchase is pending and carries an
// intent (and Payssage answered); the purchases table is the source of truth.
type snapshotPayload struct {
	PurchaseID    uuid.UUID             `json:"purchase_id"`
	EditionID     uuid.UUID             `json:"edition_id"`
	Status        models.PurchaseStatus `json:"status"`
	StatusReason  *string               `json:"status_reason,omitempty"`
	TotalCents    int64                 `json:"total_cents"`
	Currency      string                `json:"currency"`
	PaymentMethod *string               `json:"payment_method,omitempty"`
	IntentID      *uuid.UUID            `json:"intent_id,omitempty"`
	IntentStatus  *string               `json:"intent_status,omitempty"`
	Items         []snapshotItem        `json:"items"`
	QRCode        *string               `json:"qr_code,omitempty"`
	QRCodeBase64  *string               `json:"qr_code_base64,omitempty"`
	ExpiresAt     time.Time             `json:"expires_at"`
	CreatedAt     time.Time             `json:"created_at"`
}

// snapshotItem is one ledger row, same shape as the resume's items.
type snapshotItem struct {
	ItemType       models.PurchaseItemType `json:"item_type"`
	ItemID         uuid.UUID               `json:"item_id"`
	Quantity       int                     `json:"quantity"`
	UnitPriceCents int64                   `json:"unit_price_cents"`
}

// intentUpdatedPayload — pushed on intent status change only (the hub
// dedupes same-status re-deliveries). For approved → "succeeded"; for
// cancelled → "cancelled" (the specific failure vocabulary rides in
// purchase.cancelled's status_detail).
type intentUpdatedPayload struct {
	PurchaseID uuid.UUID `json:"purchase_id"`
	IntentID   uuid.UUID `json:"intent_id"`
	Status     string    `json:"status"`
}

// purchaseEventPayload — confirmed/expired/cancelled frames. The cancelled
// frame carries the provider's status_detail (insufficient_funds,
// invalid_card, high_risk…) fetched via GetIntent, best-effort.
type purchaseEventPayload struct {
	PurchaseID   uuid.UUID `json:"purchase_id"`
	StatusDetail *string   `json:"status_detail,omitempty"`
}
