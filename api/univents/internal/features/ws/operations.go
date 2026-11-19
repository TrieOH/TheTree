// Package ws is the ws feature: its handlers and operations. Build wires it.
//
// The handlers serve the realtime surfaces (split 6): the generated
// getWsToken op (REST) and the raw `WS /ws?token=...` route (mounted by the
// app router, outside the strict handler — a WebSocket handshake cannot
// carry Authorization headers and streaming must bypass the envelope
// machinery).
//
// The operations are the per-purchase WebSocket surface (split 6, D9): the only
// live-update channel for a buyer. One socket per purchase — it opens at
// checkout (the ws_token comes in the checkout response) and a disconnect is
// repaired by reconnecting to the same purchase with a fresh token (canonical
// case: a mobile user hops to the bank app to pay pix, comes back). The
// socket delivers a snapshot frame on connect (so a reconnect restores
// context), then live frames until a terminal event (confirmed/expired/
// cancelled) closes it.
//
// The socket is fed by the notifier (LISTEN/NOTIFY, kind="purchase" payloads
// published by the webhook receiver / checkout / expiry) — multi-instance
// safe: each instance's hub subscribes and fans out to its own conns.
package ws

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"lib/telemetry"
	"log/slog"
	"net/http"
	payssage "sdk/payssage"
	"sync"
	"time"
	"univents/internal/notify"
	"univents/models"
	"univents/ports"
	"uuid"

	"github.com/coder/websocket"
)

// tokenTTL is how long a handshake token lives. Short-lived because it rides
// the query string (logs): long enough for a checkout → socket open or a
// reconnect hop to the bank app and back.
const tokenTTL = 10 * time.Minute

// tokenBytes is the entropy of the raw token (32 bytes → 43 chars base64url).
const tokenBytes = 32

// Notifier is the LISTEN/NOTIFY subscribe surface (split 6 subscribes to
// kind="purchase"). Satisfied by *database.Notifier (lib/go/database).
type Notifier interface {
	Subscribe(channel string, handler func(payload string))
}

// IntentClient is the Payssage intent read seam, used for the snapshot's
// best-effort intent_status and the cancelled frame's status_detail.
// Satisfied by *payssage.Client and faked in tests.
type IntentClient interface {
	GetIntent(ctx context.Context, intentID uuid.UUID) (*payssage.Intent, error)
}

type Deps struct {
	Tokens    ports.WsTokenRepo
	Purchases ports.PurchaseRepo
	Intents   IntentClient
	// Notifier carries the store's NOTIFY channel; New subscribes the hub.
	Notifier Notifier
}

type Operations struct {
	Deps

	hub *Hub
	// intentAttempts / intentRetryDelay are the snapshot's intent-status
	// fetch budget (mirrors the split-5 resume): up to 3 tries with a short
	// pause, then degrade gracefully — a Payssage blip must not fail the
	// snapshot, the purchases table is the source of truth.
	intentAttempts   int
	intentRetryDelay time.Duration
}

func New(d Deps) *Operations {
	o := &Operations{
		Deps:             d,
		intentAttempts:   3,
		intentRetryDelay: 250 * time.Millisecond,
	}
	o.hub = newHub(d.Purchases, d.Intents)
	// One subscription for the life of the process; the handler spawns a
	// goroutine so the notifier's listen loop stays free (the fan-out may
	// hit the DB and Payssage).
	d.Notifier.Subscribe(notify.Channel, func(payload string) {
		go o.hub.handleNotification(payload)
	})
	return o
}

// SetIntentRetry overrides the snapshot intent-status retry budget. Test-only.
func (o *Operations) SetIntentRetry(attempts int, delay time.Duration) {
	o.intentAttempts = attempts
	o.intentRetryDelay = delay
}

// IssueToken is getWsToken: the caller must be the purchase's purchaser
// (anything else is NOT_FOUND via the repo — no existence leak). Issues a
// fresh one-time token — a random 32-byte value stored as its SHA-256 hash,
// valid 10 minutes — and returns the raw value (the only time it exists).
// Transaction-aware through the shared tx runner: checkout (split 7) calls
// the same method inside its tx and the token + purchase commit atomically.
func (o *Operations) IssueToken(ctx context.Context, purchaseID, userID uuid.UUID) (string, time.Time, error) {
	_, err := o.Purchases.GetByIDForOwner(ctx, purchaseID, userID)
	if err != nil {
		return "", time.Time{}, err
	}

	raw, err := newToken()
	if err != nil {
		return "", time.Time{}, err
	}
	expires := time.Now().UTC().Add(tokenTTL)
	_, err = o.Tokens.Create(ctx, &models.WsToken{
		PurchaseID: purchaseID,
		UserID:     userID,
		TokenHash:  hashToken(raw),
		ExpiresAt:  expires,
	})
	if err != nil {
		return "", time.Time{}, err
	}
	return raw, expires, nil
}

// newToken returns a random URL-safe token. crypto/rand failure is a
// genuine entropy failure — surfaced as an error, never a weak token.
func newToken() (string, error) {
	b := make([]byte, tokenBytes)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// hashToken is the at-rest form: SHA-256 hex. The raw token is never
// persisted or logged.
func hashToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

// Frame types (D9): the snapshot restores context on connect, then live
// frames until a terminal frame (confirmed/expired/cancelled) closes the
// socket. Every frame is `{"type": "...", "payload": {...}}`.
const (
	frameSnapshot  = "purchase.snapshot"
	frameIntent    = "intent.updated"
	frameConfirmed = "purchase.confirmed"
	frameExpired   = "purchase.expired"
	frameCancelled = "purchase.cancelled"
	frameFailed    = "purchase.failed"
	frameRejected  = "purchase.rejected"
)

// intent statuses the hub emits in intent.updated. Mirrors the payssage
// vocabulary for the states univents ever sees (helpers.go NormalizeStatus);
// the hub derives them from the purchase status the notification carries.
const (
	intentSucceeded = "succeeded"
	intentCancelled = "cancelled"
	intentFailed    = "failed"
	intentRejected  = "rejected"
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

const (
	// pingInterval / pongWait are the liveness budget: a ping every ~30s
	// (library-level, auto-answered by the client), and the socket is torn
	// down when a ping gets no pong within pongWait — the idle timeout.
	pingInterval = 30 * time.Second
	pongWait     = 15 * time.Second

	// writeTimeout bounds every single frame write; a stuck client must not
	// block the fan-out goroutine forever.
	writeTimeout = 10 * time.Second

	// sendBuffer is the outbound queue per conn. The notify fan-out is
	// non-blocking: a client that can't keep up is closed, never a blocker.
	sendBuffer = 32
)

// Hub owns the per-purchase sockets: purchase_id → set of conns. Events fan
// out from the notifier subscription (kind="purchase", D9); a terminal event
// (confirmed/expired/cancelled) is delivered and closes the socket. One hub
// per process; multi-instance safe because conns exist only on the instance
// the client connected to, and the notifier delivers to every instance.
type Hub struct {
	purchases ports.PurchaseRepo
	intents   IntentClient

	mu         sync.Mutex
	clients    map[uuid.UUID]map[*client]struct{}
	intentIDs  map[uuid.UUID]uuid.UUID // purchase → intent (immutable after checkout; cached at connect)
	lastStatus map[uuid.UUID]string    // dedupes intent.updated on same-status re-deliveries
}

func newHub(purchases ports.PurchaseRepo, intents IntentClient) *Hub {
	return &Hub{
		purchases:  purchases,
		intents:    intents,
		clients:    make(map[uuid.UUID]map[*client]struct{}),
		intentIDs:  make(map[uuid.UUID]uuid.UUID),
		lastStatus: make(map[uuid.UUID]string),
	}
}

// client is one socket. `send` is the outbound queue (the notify fan-out
// enqueues; the per-conn goroutine writes); `done` closes when the client
// goes away so the loop stops waiting on a dead conn.
type client struct {
	purchaseID uuid.UUID
	conn       *websocket.Conn
	send       chan outbound
	done       chan struct{}
	closeOnce  sync.Once
}

// outbound is one queued frame; terminal frames close the socket after
// being written (ordered after the frames before them).
type outbound struct {
	typ      string
	data     []byte
	terminal bool
}

func (c *client) close() {
	c.closeOnce.Do(func() {
		close(c.done)
		_ = c.conn.CloseNow()
	})
}

// register binds a conn to its purchase and caches the purchase's intent id
// (immutable after checkout) so the notify path never needs the DB for it.
func (h *Hub) register(conn *websocket.Conn, purchase *models.Purchase) *client {
	c := &client{
		purchaseID: purchase.ID,
		conn:       conn,
		send:       make(chan outbound, sendBuffer),
		done:       make(chan struct{}),
	}
	h.mu.Lock()
	if h.clients[purchase.ID] == nil {
		h.clients[purchase.ID] = make(map[*client]struct{})
	}
	h.clients[purchase.ID][c] = struct{}{}
	if purchase.PayssageIntentID != nil {
		h.intentIDs[purchase.ID] = *purchase.PayssageIntentID
	}
	h.mu.Unlock()
	return c
}

// unregister removes a conn and drops the per-purchase bookkeeping when the
// last conn leaves. Idempotent.
func (h *Hub) unregister(c *client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	set := h.clients[c.purchaseID]
	if _, ok := set[c]; !ok {
		return
	}
	delete(set, c)
	if len(set) == 0 {
		delete(h.clients, c.purchaseID)
		delete(h.intentIDs, c.purchaseID)
		delete(h.lastStatus, c.purchaseID)
	}
}

// serveClient runs the socket: background reads (control frames, disconnect
// detection), keepalive pings, and the outbound write loop. Blocks until the
// client leaves or a terminal frame closes the socket.
func (h *Hub) serveClient(c *client) {
	readCtx := c.conn.CloseRead(context.Background())
	go func() {
		<-readCtx.Done()
		c.close() // client disconnected — the write loop exits via done
	}()
	go h.keepalive(c)

	for {
		select {
		case msg := <-c.send:
			ctx, cancel := context.WithTimeout(context.Background(), writeTimeout)
			err := c.conn.Write(ctx, websocket.MessageText, msg.data)
			cancel()
			if err != nil {
				h.unregister(c)
				c.close()
				return
			}
			if msg.terminal {
				_ = c.conn.Close(websocket.StatusNormalClosure, "purchase resolved")
				h.unregister(c)
				c.close()
				return
			}
		case <-c.done:
			return
		}
	}
}

// keepalive pings every pingInterval; a ping without a pong within pongWait
// means the client is gone (or unreachable) — tear the socket down.
func (h *Hub) keepalive(c *client) {
	t := time.NewTicker(pingInterval)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			ctx, cancel := context.WithTimeout(context.Background(), pongWait)
			err := c.conn.Ping(ctx)
			cancel()
			if err != nil {
				h.unregister(c)
				c.close()
				return
			}
		case <-c.done:
			return
		}
	}
}

// handleNotification routes one notifier payload: kind="purchase" events fan
// out to the purchase's conns. Runs on its own goroutine (spawned by the
// Subscribe handler) so the listen loop is never blocked by DB/Payssage work.
func (h *Hub) handleNotification(payload string) {
	kind, _, purchase := notify.Parse(payload)
	if kind != notify.KindPurchase || purchase == nil {
		return
	}

	h.mu.Lock()
	// Snapshot the client set under the lock — enqueue iterates it without
	// the lock, so it must never be the live map: a concurrent unregister
	// (disconnect, drop) would write to a map being iterated (panic / race).
	clients := make([]*client, 0, len(h.clients[purchase.PurchaseID]))
	for c := range h.clients[purchase.PurchaseID] {
		clients = append(clients, c)
	}
	if len(clients) == 0 {
		h.mu.Unlock()
		return
	}
	intentID, hasIntent := h.intentIDs[purchase.PurchaseID]
	statusChanged := h.lastStatus[purchase.PurchaseID] != purchase.Status
	h.lastStatus[purchase.PurchaseID] = purchase.Status
	h.mu.Unlock()

	if !hasIntent {
		// Cache miss (connect raced the notification) — one PK read; the
		// intent id is immutable so this converges immediately.
		p, err := h.purchases.GetByID(context.Background(), purchase.PurchaseID)
		if err == nil && p.PayssageIntentID != nil {
			intentID, hasIntent = *p.PayssageIntentID, true
		}
	}

	frames := h.framesFor(purchase.Status, purchase.PurchaseID, intentID, hasIntent, statusChanged)
	h.enqueue(clients, frames)
}

// framesFor maps a purchase status to the frames to push (D9):
//
//	approved   → [intent.updated(succeeded), purchase.confirmed]
//	cancelled  → [intent.updated(cancelled), purchase.cancelled(+status_detail)]
//	expired    → [purchase.expired] (the intent is untouched — expiry doesn't
//	             cancel it, so no intent.updated)
//
// intent.updated fires only on status change (dedupe same-status
// re-deliveries). The last frame is terminal: the socket closes after it.
func (h *Hub) framesFor(status string, purchaseID, intentID uuid.UUID, hasIntent, statusChanged bool) []outbound {
	//nolint:exhaustive // derived from the notifier's purchase status string
	switch models.PurchaseStatus(status) {
	case models.PurchaseStatusApproved:
		frames := make([]outbound, 0, 2)
		if hasIntent && statusChanged {
			frames = append(frames, h.intentFrame(intentUpdatedPayload{purchaseID, intentID, intentSucceeded}))
		}
		frames = append(frames, h.terminalFrame(frameConfirmed, purchaseEventPayload{PurchaseID: purchaseID}))
		return frames
	case models.PurchaseStatusCancelled:
		frames := make([]outbound, 0, 2)
		if hasIntent && statusChanged {
			frames = append(frames, h.intentFrame(intentUpdatedPayload{purchaseID, intentID, intentCancelled}))
		}
		frames = append(frames, h.terminalFrame(frameCancelled, purchaseEventPayload{
			PurchaseID:   purchaseID,
			StatusDetail: h.statusDetail(intentID, hasIntent),
		}))
		return frames
	case models.PurchaseStatusFailed:
		frames := make([]outbound, 0, 2)
		if hasIntent && statusChanged {
			frames = append(frames, h.intentFrame(intentUpdatedPayload{purchaseID, intentID, intentFailed}))
		}
		frames = append(frames, h.terminalFrame(frameFailed, purchaseEventPayload{
			PurchaseID:   purchaseID,
			StatusDetail: h.statusDetail(intentID, hasIntent),
		}))
		return frames
	case models.PurchaseStatusRejected:
		frames := make([]outbound, 0, 2)
		if hasIntent && statusChanged {
			frames = append(frames, h.intentFrame(intentUpdatedPayload{purchaseID, intentID, intentRejected}))
		}
		frames = append(frames, h.terminalFrame(frameRejected, purchaseEventPayload{
			PurchaseID:   purchaseID,
			StatusDetail: h.statusDetail(intentID, hasIntent),
		}))
		return frames
	case models.PurchaseStatusExpired:
		return []outbound{h.terminalFrame(frameExpired, purchaseEventPayload{PurchaseID: purchaseID})}
	default:
		telemetry.Log().Warn("ws: unhandled purchase status in notification",
			slog.String("status", status),
			slog.String("purchase_id", purchaseID.String()))
		return nil
	}
}

// statusDetail fetches the provider's failure vocabulary for the cancelled
// frame (GetIntent → status_detail). Best-effort with a bounded budget: a
// Payssage blip must not stall the fan-out — omit the detail, the frame
// still goes out.
func (h *Hub) statusDetail(intentID uuid.UUID, hasIntent bool) *string {
	if !hasIntent {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	intent, err := h.intents.GetIntent(ctx, intentID)
	if err != nil {
		telemetry.Log().Warn("ws: intent status_detail unavailable", telemetry.Err(err))
		return nil
	}
	return intent.StatusDetail
}

// enqueue fans the frames out to every conn. Non-blocking: a full queue
// means the client isn't keeping up — drop it (the snapshot-on-reconnect
// restores context) and stop queueing the remaining frames for it.
func (h *Hub) enqueue(clients []*client, frames []outbound) {
	for _, c := range clients {
	forFrame:
		for _, f := range frames {
			select {
			case c.send <- f:
			default:
				h.unregister(c)
				c.close()
				break forFrame
			}
		}
	}
}

func (h *Hub) intentFrame(payload intentUpdatedPayload) outbound {
	return outbound{typ: frameIntent, data: mustFrame(frameIntent, payload)}
}

func (h *Hub) terminalFrame(typ string, payload purchaseEventPayload) outbound {
	return outbound{typ: typ, data: mustFrame(typ, payload), terminal: true}
}

// mustFrame marshals a frame; the payload structs are static shapes so
// marshaling cannot fail — a failure here would be a programming error, so
// it surfaces loudly instead of silently dropping a frame.
func mustFrame(typ string, payload any) []byte {
	b, err := marshalFrame(typ, payload)
	if err != nil {
		slog.Error("ws: marshal frame", "type", typ, "err", err)
		panic(err)
	}
	return b
}

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

	tok, err := o.Tokens.Consume(r.Context(), hashToken(raw))
	if err != nil {
		telemetry.Log().Error("ws: token consume failed",
			slog.String("purchase_id", r.URL.Query().Get("purchase_id")),
			telemetry.Err(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	if tok == nil {
		// One-time by construction: no row means missing, already used, or
		// expired — all reject. The client re-issues via GET /ws/token.
		http.Error(w, "invalid or expired token", http.StatusUnauthorized)
		return
	}

	purchase, err := o.Purchases.GetByID(r.Context(), tok.PurchaseID)
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
			slog.String("purchase_id", purchase.ID.String()),
			telemetry.Err(err))
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
	items, err := o.Purchases.ListItemsByPurchase(ctx, purchase.ID)
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
		intent, err := o.Intents.GetIntent(ctx, intentID)
		if err == nil {
			return new(string(intent.Status))
		}
		lastErr = err
		if attempt < o.intentAttempts {
			select {
			case <-ctx.Done():
				telemetry.Log().Warn("ws: intent status fetch cancelled",
					slog.String("intent_id", intentID.String()),
					telemetry.Err(ctx.Err()))
				return nil
			case <-time.After(o.intentRetryDelay):
			}
		}
	}
	telemetry.Log().Error("ws: intent status unavailable after retries — degrading (purchases table remains the source of truth)",
		slog.String("intent_id", intentID.String()),
		slog.Int("attempts", o.intentAttempts),
		telemetry.Err(lastErr))
	return nil
}
