// Package store is the storefront's live-stock surface (split 6, D10): one
// SSE stream per edition, public. On connect it streams a full snapshot of
// every purchasable item with its computed `{id, item_type, stock}`
// (availability semantics from split 3), then relays stock deltas published
// by the notifier (kind="stock") — one relay per edition, shared by its
// conns, re-querying availability from the DB so the numbers are always
// authoritative (publishers carry item_ids only, never stock numbers).
package store

import (
	"context"
	"encoding/json"
	"net/http"

	"lib/telemetry"

	"univents/models"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Notifier is the LISTEN/NOTIFY subscribe surface (split 6 subscribes to
// kind="stock"). Satisfied by *database.Notifier (lib/go/database).
type Notifier interface {
	Subscribe(channel string, handler func(payload string))
}

// AvailabilityReader is the stock read the SSE surface needs — the full
// purchase ledger, filtered to what split 6 touches. Satisfied by
// *repos/purchases.Repo (ports.PurchaseRepo).
type AvailabilityReader interface {
	Availability(ctx context.Context, editionID uuid.UUID) ([]models.ItemAvailability, error)
}

// EditionRepo resolves the route's edition_id to a 404 when unknown.
type EditionRepo interface {
	GetByID(ctx context.Context, id uuid.UUID) (*models.Edition, error)
}

type Operations struct {
	avail    AvailabilityReader
	editions EditionRepo
	manager  *RelayManager
}

func NewOperations(avail AvailabilityReader, editions EditionRepo, notifier Notifier) *Operations {
	o := &Operations{
		avail:    avail,
		editions: editions,
	}
	o.manager = newRelayManager(avail, notifier)
	return o
}

// Stock is the REST twin of the SSE snapshot — every purchasable item's
// current stock position, recomputed from the availability ledger (base −
// reserved across pending/approved purchases; null = unlimited). Public:
// the storefront needs the counts before/without subscribing to the stream.
// Unknown editions are NOT_FOUND (the endpoint exists, the edition does not).
func (o *Operations) Stock(ctx context.Context, editionID uuid.UUID) ([]stockItem, error) {
	ctx, span := telemetry.StartSpan(ctx, "StoreService.Stock")
	defer span.End()

	if _, err := o.editions.GetByID(ctx, editionID); err != nil {
		return nil, err
	}
	avail, err := o.avail.Availability(ctx, editionID)
	if err != nil {
		return nil, err
	}
	items := make([]stockItem, 0, len(avail))
	for _, a := range avail {
		items = append(items, stockItem{
			ID:       a.ItemID,
			ItemType: string(a.ItemType),
			Stock:    available(a),
		})
	}
	return items, nil
}

// ServeStream is the raw SSE route — `GET /editions/{edition_id}/store/stream`
// (public). Raw on purpose: a streaming response bypasses the
// fun/validate-envelope machinery (which would buffer the body), and the
// harness's WriteTimeout would kill a long-lived stream — the handler
// hijacks the connection (clearing all deadlines) and speaks close-delimited
// HTTP/1.1 chunk-free SSE directly.
func (o *Operations) ServeStream(w http.ResponseWriter, r *http.Request, editionID uuid.UUID) {
	_, err := o.editions.GetByID(r.Context(), editionID)
	if err != nil {
		http.Error(w, "edition not found", http.StatusNotFound)
		return
	}

	stream, err := newSSEStream(w)
	if err != nil {
		telemetry.Log().Error("store: hijack failed", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	relay := o.manager.relayFor(editionID)
	sub := relay.subscribe(stream)
	defer sub.close()

	// Snapshot first (re-read from the DB — the same query the relay uses,
	// so the client's first frame is consistent with the deltas that follow).
	err = o.writeSnapshot(r.Context(), stream, editionID)
	if err != nil {
		sub.close()
		return
	}

	// Keepalive comment line (": ping") every ~30s and disconnect detection
	// on the hijacked conn.
	go stream.keepalive()
	stream.waitClosed()
}

// writeSnapshot streams `event: snapshot` with every purchasable item's
// current stock position (available = base − reserved; null = unlimited).
func (o *Operations) writeSnapshot(ctx context.Context, stream *sseStream, editionID uuid.UUID) error {
	avail, err := o.avail.Availability(ctx, editionID)
	if err != nil {
		telemetry.Log().Error("store: snapshot availability read failed",
			zap.String("edition_id", editionID.String()),
			zap.Error(err))
		return err
	}

	items := make([]stockItem, 0, len(avail))
	for _, a := range avail {
		items = append(items, stockItem{
			ID:       a.ItemID,
			ItemType: string(a.ItemType),
			Stock:    available(a),
		})
	}
	raw, err := json.Marshal(items)
	if err != nil {
		return err
	}
	return stream.event("snapshot", raw)
}

// stockItem is one SSE stock payload (decision at review): the item's id
// plus its type tag — ids span three tables (ticket_types /
// product_variants / program_occurrences) and the front must know which
// catalog to match against without a join.
type stockItem struct {
	ID       uuid.UUID `json:"id"`
	ItemType string    `json:"item_type"`
	Stock    *int      `json:"stock"` // null = unlimited
}

// available computes available = base − reserved; nil base = unlimited.
func available(a models.ItemAvailability) *int {
	if a.BaseQuantity == nil {
		return nil
	}
	avail := *a.BaseQuantity - int(a.ReservedQuantity)
	return &avail
}
