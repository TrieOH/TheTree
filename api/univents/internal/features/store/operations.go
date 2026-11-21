// Package store is the store feature: its handlers and operations. Build
// wires it.
//
// The handlers serve the storefront's live-stock SSE surface (split 6):
// the raw `GET /editions/{edition_id}/store/stream` route (mounted by the
// app router, outside the strict handler — streaming must bypass the
// fun/validate envelope machinery).
//
// The operations are the storefront's live-stock surface (split 6, D10): one
// SSE stream per edition, public. On connect it streams a full snapshot of
// every purchasable item with its computed `{id, item_type, stock}`
// (availability semantics from split 3), then relays stock deltas published
// by the notifier (kind="stock") — one relay per edition, shared by its
// conns, re-querying availability from the DB so the numbers are always
// authoritative (publishers carry item_ids only, never stock numbers).
package store

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"lib/telemetry"
	"log/slog"
	"net"
	"net/http"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"univents/internal/notify"
	"univents/models"
	"uuid"
)

// Notifier is the LISTEN/NOTIFY subscribe surface (split 6 subscribes to
// kind="stock"). Satisfied by *database.Notifier (lib/go/database).
type Notifier interface {
	Subscribe(channel string, handler func(payload string))
}

// AvailabilityReader is the stock read the SSE surface needs — the full
// purchase ledger, filtered to what split 6 touches. Satisfied by
// *purchases.Repo (ports.PurchaseRepo).
type AvailabilityReader interface {
	Availability(ctx context.Context, editionID uuid.UUID) ([]models.ItemAvailability, error)
}

// EditionRepo resolves the route's edition_id to a 404 when unknown.
type EditionRepo interface {
	GetByID(ctx context.Context, id uuid.UUID) (*models.Edition, error)
}

type Deps struct {
	Avail    AvailabilityReader
	Editions EditionRepo
	// Notifier carries the store's NOTIFY channel; New hands it to the relay.
	Notifier Notifier
}

type Operations struct {
	Deps

	manager *RelayManager
}

func New(d Deps) *Operations {
	return &Operations{Deps: d, manager: newRelayManager(d.Avail, d.Notifier)}
}

// Stock is the REST twin of the SSE snapshot — every purchasable item's
// current stock position, recomputed from the availability ledger (base −
// reserved across pending/approved purchases; null = unlimited). Public:
// the storefront needs the counts before/without subscribing to the stream.
// Unknown editions are NOT_FOUND (the endpoint exists, the edition does not).
func (o *Operations) Stock(ctx context.Context, editionID uuid.UUID) ([]stockItem, error) {
	_, err := o.Editions.GetByID(ctx, editionID)
	if err != nil {
		return nil, err
	}
	avail, err := o.Avail.Availability(ctx, editionID)
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
	_, err := o.Editions.GetByID(r.Context(), editionID)
	if err != nil {
		http.Error(w, "edition not found", http.StatusNotFound)
		return
	}

	stream, err := newSSEStream(w)
	if err != nil {
		telemetry.Log().ErrorContext(r.Context(), "store: hijack failed", telemetry.Err(err))
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
	avail, err := o.Avail.Availability(ctx, editionID)
	if err != nil {
		telemetry.Log().ErrorContext(ctx, "store: snapshot availability read failed",
			slog.String("edition_id", editionID.String()),
			telemetry.Err(err))
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

// RelayManager owns the stock-delta fan-out: ONE notifier subscription for
// kind="stock" (registered at construction, never torn down) routing to one
// Relay per edition. Relays are created lazily on the first SSE conn and
// dropped when the last conn leaves — the subscription itself lives for the
// process, so there is nothing to unsubscribe (the notifier has no per-
// handler teardown).
type RelayManager struct {
	avail    AvailabilityReader
	notifier Notifier

	mu     sync.Mutex
	relays map[uuid.UUID]*Relay
}

func newRelayManager(avail AvailabilityReader, notifier Notifier) *RelayManager {
	m := &RelayManager{
		avail:    avail,
		notifier: notifier,
		relays:   make(map[uuid.UUID]*Relay),
	}
	notifier.Subscribe(notify.Channel, func(payload string) {
		// Off the listen loop: the relay re-queries the DB, and a missed or
		// reordered delta is harmless — every broadcast re-reads the current
		// availability, so the numbers are always the authoritative truth.
		go m.handleStock(payload)
	})
	return m
}

// relayFor returns the edition's shared relay, creating it on first use.
func (m *RelayManager) relayFor(editionID uuid.UUID) *Relay {
	m.mu.Lock()
	defer m.mu.Unlock()
	if r, ok := m.relays[editionID]; ok {
		return r
	}
	r := &Relay{editionID: editionID, avail: m.avail, manager: m}
	m.relays[editionID] = r
	return r
}

func (m *RelayManager) handleStock(payload string) {
	kind, stock, _ := notify.Parse(payload)
	if kind != notify.KindStock || stock == nil || len(stock.ItemIDs) == 0 {
		return
	}

	m.mu.Lock()
	relay, ok := m.relays[stock.EditionID]
	m.mu.Unlock()
	if !ok {
		return // no conns for this edition — nothing to update
	}
	relay.broadcast(stock)
}

func (m *RelayManager) drop(relay *Relay) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if r, ok := m.relays[relay.editionID]; ok && r == relay {
		delete(m.relays, relay.editionID)
	}
}

// Relay is one edition's shared delta fan-out: its conns all get the same
// broadcast, and the availability re-query happens once per NOTIFY
// regardless of conn count (D10).
type Relay struct {
	editionID uuid.UUID
	avail     AvailabilityReader
	manager   *RelayManager

	mu    sync.Mutex
	conns map[*sseStream]struct{}
}

// subscribe registers a stream; the returned handle unregisters on close
// and drops the relay from the manager when the last conn leaves.
func (r *Relay) subscribe(s *sseStream) *subscription {
	r.mu.Lock()
	if r.conns == nil {
		r.conns = make(map[*sseStream]struct{})
	}
	r.conns[s] = struct{}{}
	r.mu.Unlock()
	return &subscription{relay: r, stream: s}
}

type subscription struct {
	relay  *Relay
	stream *sseStream
}

func (s *subscription) close() {
	r := s.relay
	r.mu.Lock()
	delete(r.conns, s.stream)
	empty := len(r.conns) == 0
	r.mu.Unlock()
	if empty {
		r.manager.drop(r)
	}
}

// broadcast re-queries availability for the NOTIFY'd item_ids and writes one
// `event: stock` per item to every conn. The DB read is the authority: the
// publisher's item_ids are a hint, the numbers come from here (D10).
func (r *Relay) broadcast(stock *notify.StockNotification) {
	availability, err := r.avail.Availability(context.Background(), r.editionID)
	if err != nil {
		telemetry.Log().ErrorContext(context.Background(), "store: relay availability read failed",
			slog.String("edition_id", r.editionID.String()),
			telemetry.Err(err))
		return
	}

	// Filter the full availability read down to the affected items. An id
	// that misses (wrong edition, deleted item) is skipped — its stock
	// simply does not update.
	affected := make(map[uuid.UUID]models.ItemAvailability, len(stock.ItemIDs))
	for _, a := range availability {
		if slices.Contains(stock.ItemIDs, a.ItemID) {
			affected[a.ItemID] = a
		}
	}
	if len(affected) == 0 {
		return
	}

	r.mu.Lock()
	conns := make([]*sseStream, 0, len(r.conns))
	for c := range r.conns {
		conns = append(conns, c)
	}
	r.mu.Unlock()

	for _, a := range affected {
		raw, err := json.Marshal(stockItem{ID: a.ItemID, ItemType: string(a.ItemType), Stock: available(a)})
		if err != nil {
			continue
		}
		for _, c := range conns {
			// Best-effort: a dead conn drops the write and is removed by
			// waitClosed — never block the relay on a stuck stream.
			_ = c.event("stock", raw)
		}
	}
}

// sseKeepaliveIntervalNS is how often the stream emits a `: ping` comment
// line. Atomic so tests can shorten it without racing in-flight streams
// (a previous test's stream goroutine may still be reading it).
var sseKeepaliveIntervalNS atomic.Int64

func init() { sseKeepaliveIntervalNS.Store(int64(30 * time.Second)) }

// sseStream is a hijacked close-delimited HTTP/1.1 response speaking SSE.
// The harness's middleware stack wraps the ResponseWriter in writers that
// expose no Flush and no write-deadline control, and its http.Server sets a
// 60s WriteTimeout — all fatal for a long-lived stream. Hijacking the conn
// clears every deadline (net/http does that on Hijack) and hands the raw
// connection to us: headers are written once, then each event is a framed
// write + flush. No Content-Length / Transfer-Encoding means the body is
// close-delimited — legal HTTP/1.1 and exactly what SSE wants.
//
// Only one goroutine may write at a time; event() takes the mutex.
type sseStream struct {
	conn net.Conn
	bw   *bufio.Writer
	mu   sync.Mutex

	done      chan struct{}
	closeOnce sync.Once
}

// corsResponseHeaders are the response headers the harness CORS middleware
// (rs/cors with AllowCredentials) sets on the ResponseWriter for an actual
// request. Hijack discards the ResponseWriter, so the hand-built head must
// re-emit them or the browser blocks the cross-origin stream (EventSource
// is CORS-gated). Connection/Content-Length/Transfer-Encoding are hop-by-hop
// and deliberately excluded: the body is close-delimited.
var corsResponseHeaders = []string{
	"Access-Control-Allow-Origin",
	"Access-Control-Allow-Credentials",
	"Access-Control-Expose-Headers",
	"Vary",
}

// newSSEStream hijacks the connection and writes the response head. Any
// error means the connection was not taken over (the caller can still write
// a plain HTTP error).
func newSSEStream(w http.ResponseWriter) (*sseStream, error) {
	hj, ok := w.(http.Hijacker)
	if !ok {
		return nil, errors.New("response writer does not support hijacking")
	}
	conn, rw, err := hj.Hijack()
	if err != nil {
		return nil, err
	}
	bw := rw.Writer

	// Snapshot the CORS headers the middleware set before hijacking (the
	// ResponseWriter is discarded on Hijack) and splice them into the raw
	// head; no Origin header means none were set and none are emitted.
	head := strings.Builder{}
	head.WriteString("HTTP/1.1 200 OK\r\n")
	for _, name := range corsResponseHeaders {
		for _, v := range w.Header().Values(name) {
			head.WriteString(name + ": " + v + "\r\n")
		}
	}
	head.WriteString("Content-Type: text/event-stream; charset=utf-8\r\n" +
		"Cache-Control: no-store\r\n" +
		"Connection: keep-alive\r\n" +
		"X-Accel-Buffering: no\r\n" +
		"\r\n")
	_, err = bw.WriteString(head.String())
	if err != nil {
		_ = conn.Close()
		return nil, err
	}
	err = bw.Flush()
	if err != nil {
		_ = conn.Close()
		return nil, err
	}
	return &sseStream{conn: conn, bw: bw, done: make(chan struct{})}, nil
}

// event writes one SSE event block (`event: <name>\ndata: <json>\n\n`) and
// flushes it to the wire. Returns the write error (the caller treats a dead
// conn as terminal).
func (s *sseStream) event(name string, data []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := fmt.Fprintf(s.bw, "event: %s\n", name)
	if err != nil {
		return err
	}
	// SSE forbids bare newlines inside data; JSON here is single-line, but
	// write defensively so a multi-line payload cannot desync the protocol.
	for _, line := range splitLines(data) {
		_, err := fmt.Fprintf(s.bw, "data: %s\n", line)
		if err != nil {
			return err
		}
	}
	_, err = s.bw.WriteString("\n")
	if err != nil {
		return err
	}
	return s.bw.Flush()
}

// keepalive writes a comment line (`: ping`) every ~30s so proxies and the
// client can tell the stream is alive. Exits when the stream closes.
func (s *sseStream) keepalive() {
	t := time.NewTicker(time.Duration(sseKeepaliveIntervalNS.Load()))
	defer t.Stop()
	for {
		select {
		case <-t.C:
			s.mu.Lock()
			_, _ = s.bw.WriteString(": ping\n\n")
			err := s.bw.Flush()
			s.mu.Unlock()
			if err != nil {
				s.close() // conn died under us — tear down
				return
			}
		case <-s.done:
			return
		}
	}
}

// waitClosed blocks until the client goes away (the client never sends
// data, so the first read fails only on disconnect), then closes the
// stream — the close-delimited body terminator.
func (s *sseStream) waitClosed() {
	buf := make([]byte, 1)
	_, _ = s.conn.Read(buf)
	s.close()
}

// close tears the connection down (idempotent).
func (s *sseStream) close() {
	s.closeOnce.Do(func() {
		close(s.done)
		_ = s.conn.Close()
	})
}

// splitLines splits data into lines (SSE forbids bare newlines in data).
func splitLines(b []byte) [][]byte {
	var lines [][]byte
	start := 0
	for i, c := range b {
		if c == '\n' {
			lines = append(lines, b[start:i])
			start = i + 1
		}
	}
	return append(lines, b[start:])
}
