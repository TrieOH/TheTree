package ws

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"lib/telemetry"

	"univents/internal/services/notify"
	"univents/models"
	"univents/ports"

	"github.com/coder/websocket"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

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
			zap.String("status", status),
			zap.String("purchase_id", purchaseID.String()))
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
		telemetry.Log().Warn("ws: intent status_detail unavailable", zap.Error(err))
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
