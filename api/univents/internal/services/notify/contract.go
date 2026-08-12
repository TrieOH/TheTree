// Package notify is the wire contract of the store's LISTEN/NOTIFY channel
// (`univents_changes`): the payloads publishers (webhook receiver, checkout,
// expiry) write and the realtime subscribers (WS hub, SSE relay, split 6)
// parse. The publisher and subscribers never share the Go types — the JSON
// shape is the contract — but this package keeps the shape in one place so
// the two sides cannot drift. Publishers marshal these types; subscribers
// unmarshal into them and route on `Kind`.
package notify

import (
	"encoding/json"

	"github.com/google/uuid"
)

// Channel is the single NOTIFY channel the store publishes on. The WS hub
// routes kind="purchase" events (D9) and the SSE relay routes kind="stock"
// deltas (D10), both subscribed in split 6.
const Channel = "univents_changes"

// Kinds — every payload carries exactly one.
const (
	KindStock    = "stock"
	KindPurchase = "purchase"
)

// StockNotification is the stock-delta payload (D10): item_ids only, never
// stock numbers — the SSE relay re-queries availability from the DB, so a
// missed notification is a stale snapshot, never data loss. Note: item_ids
// are UUIDs across three tables (ticket_types / product_variants /
// program_occurrences) with no type tag — the relay's re-query must look up
// across all three (split 3 availability), or an id misses its table and
// stock silently stops updating.
type StockNotification struct {
	Kind      string      `json:"kind"`
	EditionID uuid.UUID   `json:"edition_id"`
	ItemIDs   []uuid.UUID `json:"item_ids"`
}

// PurchaseNotification is the purchase-event payload (D9). The WS hub maps
// status to the frame type: approved → purchase.confirmed, expired →
// purchase.expired, cancelled → purchase.cancelled.
type PurchaseNotification struct {
	Kind       string    `json:"kind"`
	EditionID  uuid.UUID `json:"edition_id"`
	PurchaseID uuid.UUID `json:"purchase_id"`
	Status     string    `json:"status"`
}

// Parse decodes a channel payload into its kind-specific shape. Returns nil
// for payloads this channel does not route (unknown kind or malformed JSON
// — subscribers log and skip rather than crash the listen loop).
func Parse(payload string) (kind string, stock *StockNotification, purchase *PurchaseNotification) {
	raw := struct {
		Kind string `json:"kind"`
	}{}
	err := json.Unmarshal([]byte(payload), &raw)
	if err != nil {
		return "", nil, nil
	}
	switch raw.Kind {
	case KindStock:
		var n StockNotification
		err = json.Unmarshal([]byte(payload), &n)
		if err != nil {
			return "", nil, nil
		}
		return KindStock, &n, nil
	case KindPurchase:
		var n PurchaseNotification
		err = json.Unmarshal([]byte(payload), &n)
		if err != nil {
			return "", nil, nil
		}
		return KindPurchase, nil, &n
	default:
		return raw.Kind, nil, nil
	}
}
