package ws

import (
	"context"
	"net/http"
	"time"

	"lib/telemetry"

	"univents/models"

	"github.com/coder/websocket"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// ServeWS is the raw WS route — `WS /ws?token=...` (D9). The token is the
// handshake auth (WS handshakes cannot carry Authorization headers): it
// proves prior REST auth for this purchase and nothing more. Handshake:
// consume the one-time token atomically (missing/already-used/expired →
// reject), bind the socket to the purchase, send the snapshot frame (so a
// reconnect restores context — no resume polling), then relay live frames
// until a terminal event closes the socket.
func (o *Operations) ServeWS(w http.ResponseWriter, r *http.Request) {
	raw := r.URL.Query().Get("token")
	if raw == "" {
		http.Error(w, "missing token", http.StatusBadRequest)
		return
	}

	tok, err := o.tokens.Consume(r.Context(), hashToken(raw))
	if err != nil {
		telemetry.Log().Error("ws: token consume failed",
			zap.String("purchase_id", r.URL.Query().Get("purchase_id")),
			zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	if tok == nil {
		// One-time by construction: no row means missing, already used, or
		// expired — all reject. The client re-issues via GET /ws/token.
		http.Error(w, "invalid or expired token", http.StatusUnauthorized)
		return
	}

	purchase, err := o.purchases.GetByID(r.Context(), tok.PurchaseID)
	if err != nil {
		http.Error(w, "purchase not found", http.StatusNotFound)
		return
	}

	// No origin check (decision at review): the one-time token is the auth.
	// Browsers send an Origin header and the CORS middleware already passed
	// the request; non-browser clients (mobile) send none.
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{InsecureSkipVerify: true})
	if err != nil {
		return // Accept wrote the response
	}

	c := o.hub.register(conn, purchase)
	defer o.hub.unregister(c)

	// Register before the snapshot: any live event landing while the
	// snapshot is built is queued and delivered after it, so the client
	// never misses a terminal event behind a stale snapshot.
	snap, err := o.snapshot(r.Context(), purchase)
	if err != nil {
		telemetry.Log().Error("ws: snapshot build failed",
			zap.String("purchase_id", purchase.ID.String()),
			zap.Error(err))
		_ = conn.Close(websocket.StatusInternalError, "snapshot failed")
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), writeTimeout)
	err = conn.Write(ctx, websocket.MessageText, snap)
	cancel()
	if err != nil {
		return
	}

	o.hub.serveClient(c)
}

// snapshot builds the connect frame: the purchase's current state (one row
// read + items) and — only when pending and carrying an intent — the
// intent's status, best-effort (same graceful-degradation budget as the
// split-5 resume: a Payssage blip must not fail the snapshot).
func (o *Operations) snapshot(ctx context.Context, purchase *models.Purchase) ([]byte, error) {
	items, err := o.purchases.ListItemsByPurchase(ctx, purchase.ID)
	if err != nil {
		return nil, err
	}

	payload := snapshotPayload{
		PurchaseID:    purchase.ID,
		EditionID:     purchase.EditionID,
		Status:        purchase.Status,
		StatusReason:  purchase.StatusReason,
		TotalCents:    purchase.TotalCents,
		Currency:      purchase.Currency,
		PaymentMethod: purchase.PaymentMethod,
		IntentID:      purchase.PayssageIntentID,
		Items:         make([]snapshotItem, 0, len(items)),
		QRCode:        purchase.QRCode,
		QRCodeBase64:  purchase.QRCodeBase64,
		ExpiresAt:     purchase.ExpiresAt,
		CreatedAt:     purchase.CreatedAt,
	}
	for _, item := range items {
		payload.Items = append(payload.Items, snapshotItem{
			ItemType:       item.ItemType,
			ItemID:         item.ItemID,
			Quantity:       item.Quantity,
			UnitPriceCents: item.UnitPriceCents,
		})
	}
	if purchase.Status == models.PurchaseStatusPending && purchase.PayssageIntentID != nil {
		payload.IntentStatus = o.intentStatus(ctx, *purchase.PayssageIntentID)
	}
	return marshalFrame(frameSnapshot, payload)
}

// intentStatus fetches the intent's status with the retry budget in
// Operations (3 tries by default), then degrades gracefully: on persistent
// failure it logs loudly and returns nil — the snapshot still carries the
// purchase's own state, which is the source of truth.
func (o *Operations) intentStatus(ctx context.Context, intentID uuid.UUID) *string {
	var lastErr error
	for attempt := 1; attempt <= o.intentAttempts; attempt++ {
		intent, err := o.intents.GetIntent(ctx, intentID)
		if err == nil {
			return new(string(intent.Status))
		}
		lastErr = err
		if attempt < o.intentAttempts {
			select {
			case <-ctx.Done():
				telemetry.Log().Warn("ws: intent status fetch cancelled",
					zap.String("intent_id", intentID.String()),
					zap.Error(ctx.Err()))
				return nil
			case <-time.After(o.intentRetryDelay):
			}
		}
	}
	telemetry.Log().Error("ws: intent status unavailable after retries — degrading (purchases table remains the source of truth)",
		zap.String("intent_id", intentID.String()),
		zap.Int("attempts", o.intentAttempts),
		zap.Error(lastErr))
	return nil
}
