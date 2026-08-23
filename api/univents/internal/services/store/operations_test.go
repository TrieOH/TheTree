package store_test

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"univents/internal/services/notify"
	"univents/internal/services/store"
	"univents/models"

	"github.com/MintzyG/fun"
	mws "github.com/MintzyG/fun/middlewares"
	"github.com/google/uuid"
)

func TestMain(m *testing.M) {
	fun.SetConfig(fun.Config{
		DefaultModule:        "test",
		DefaultContentType:   "application/json",
		EnableSizeValidation: true,
	})
	os.Exit(m.Run())
}

// ── fakes ────────────────────────────────────────────────────────────────

type fakeAvailability struct {
	mu    sync.Mutex
	items []models.ItemAvailability
}

func (f *fakeAvailability) Availability(context.Context, uuid.UUID) ([]models.ItemAvailability, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]models.ItemAvailability(nil), f.items...), nil
}

func (f *fakeAvailability) set(items []models.ItemAvailability) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.items = items
}

type fakeEditions struct{ editionID uuid.UUID }

func (f *fakeEditions) GetByID(_ context.Context, id uuid.UUID) (*models.Edition, error) {
	if id != f.editionID {
		return nil, fun.ErrNotFound("edition not found")
	}
	return &models.Edition{ID: id}, nil
}

type fakeNotifier struct {
	mu       sync.Mutex
	handlers []func(payload string)
}

func (f *fakeNotifier) Subscribe(_ string, handler func(payload string)) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.handlers = append(f.handlers, handler)
}

func (f *fakeNotifier) fire(payload string) {
	f.mu.Lock()
	handlers := append([]func(string){}, f.handlers...)
	f.mu.Unlock()
	for _, h := range handlers {
		h(payload)
	}
}

// ── raw TCP SSE client ───────────────────────────────────────────────────

// sseClient speaks raw HTTP/1.1 to the hijacked stream: write the GET, read
// the response head, then read SSE event blocks line by line.
type sseClient struct {
	conn net.Conn
	rd   *bufio.Reader
}

func dialSSE(t *testing.T, ts *httptest.Server, editionID uuid.UUID) *sseClient {
	t.Helper()
	c, _ := dialSSEHead(t, ts, editionID, "")
	return c
}

// dialSSEHead is dialSSE plus the raw response head lines (status line and
// headers), so tests can assert what the hijacked stream actually emits.
// extraHead are raw request header lines to send (e.g. "Origin: …\r\n").
func dialSSEHead(t *testing.T, ts *httptest.Server, editionID uuid.UUID, extraHead string) (*sseClient, []string) {
	t.Helper()
	dialer := net.Dialer{}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	conn, err := dialer.DialContext(ctx, "tcp", ts.Listener.Addr().String())
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	t.Cleanup(func() { _ = conn.Close() })

	_, err = fmt.Fprintf(conn, "GET /editions/%s/store/stream HTTP/1.1\r\nHost: test\r\n%s\r\n", editionID, extraHead)
	if err != nil {
		t.Fatalf("write request: %v", err)
	}
	rd := bufio.NewReader(conn)
	// Response head up to the blank line.
	var head []string
	for {
		line, err := rd.ReadString('\n')
		if err != nil {
			t.Fatalf("read head: %v", err)
		}
		if line == "\r\n" {
			break
		}
		head = append(head, strings.TrimRight(line, "\r\n"))
	}
	return &sseClient{conn: conn, rd: rd}, head
}

// event is one SSE block: `event: <name>` / `data: <json>` / blank line.
type event struct {
	name string
	data string
}

// readEvent parses the next SSE block. The keepalive comment line (`: ping`)
// is skipped. Blocks until data arrives or the conn dies.
func (c *sseClient) readEvent(t *testing.T) (event, bool) {
	t.Helper()
	_ = c.conn.SetReadDeadline(time.Now().Add(5 * time.Second))

	var ev event
	var inEvent bool
	for {
		line, err := c.rd.ReadString('\n')
		if err != nil {
			return ev, false
		}
		line = strings.TrimRight(line, "\r\n")
		switch {
		case line == "":
			if inEvent {
				return ev, true
			}
		case strings.HasPrefix(line, "event: "):
			ev.name = strings.TrimPrefix(line, "event: ")
			inEvent = true
		case strings.HasPrefix(line, "data: "):
			ev.data = strings.TrimPrefix(line, "data: ")
			inEvent = true
		}
	}
}

// ── tests ────────────────────────────────────────────────────────────────

func newStore(t *testing.T, editionID uuid.UUID) (*store.Operations, *fakeAvailability, *fakeNotifier) {
	t.Helper()
	avail := &fakeAvailability{}
	notifier := &fakeNotifier{}
	ops := store.NewOperations(avail, &fakeEditions{editionID: editionID}, notifier)
	return ops, avail, notifier
}

func TestStock(t *testing.T) {
	editionID := uuid.New()
	ops, avail, _ := newStore(t, editionID)
	ticketID, variantID := uuid.New(), uuid.New()
	avail.set([]models.ItemAvailability{
		{ItemType: models.PurchaseItemTypeTicket, ItemID: ticketID, BaseQuantity: new(int(10)), ReservedQuantity: 3},
		{ItemType: models.PurchaseItemTypeProduct, ItemID: variantID, BaseQuantity: nil, ReservedQuantity: 0}, // unlimited
	})

	items, err := ops.Stock(context.Background(), editionID)
	if err != nil {
		t.Fatalf("Stock: %v", err)
	}
	if len(items) != 2 {
		t.Fatalf("items = %d, want 2", len(items))
	}
	byID := map[uuid.UUID]struct {
		ItemType string
		Stock    *int
	}{}
	for _, it := range items {
		byID[it.ID] = struct {
			ItemType string
			Stock    *int
		}{it.ItemType, it.Stock}
	}
	if s := byID[ticketID]; s.ItemType != "ticket" || s.Stock == nil || *s.Stock != 7 {
		t.Fatalf("ticket = %+v, want ticket stock 7 (10-3)", s)
	}
	if s := byID[variantID]; s.ItemType != "product" || s.Stock != nil {
		t.Fatalf("unlimited variant = %+v, want product stock null", s)
	}
}

func TestStockUnknownEdition(t *testing.T) {
	ops, _, _ := newStore(t, uuid.New())
	_, err := ops.Stock(context.Background(), uuid.New())
	if err == nil {
		t.Fatal("Stock on unknown edition: want error, got nil")
	}
}

func TestSSESnapshotShape(t *testing.T) {
	editionID := uuid.New()
	ops, avail, _ := newStore(t, editionID)
	ticketID, variantID, occurrenceID := uuid.New(), uuid.New(), uuid.New()
	avail.set([]models.ItemAvailability{
		{ItemType: models.PurchaseItemTypeTicket, ItemID: ticketID, BaseQuantity: new(int(10)), ReservedQuantity: 3},
		{ItemType: models.PurchaseItemTypeProduct, ItemID: variantID, BaseQuantity: nil, ReservedQuantity: 0}, // unlimited
		{ItemType: models.PurchaseItemTypeProgramOccurrence, ItemID: occurrenceID, BaseQuantity: new(int(5)), ReservedQuantity: 0},
	})

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ops.ServeStream(w, r, editionID)
	}))
	t.Cleanup(ts.Close)

	c := dialSSE(t, ts, editionID)
	ev, ok := c.readEvent(t)
	if !ok {
		t.Fatal("no snapshot event")
	}
	if ev.name != "snapshot" {
		t.Fatalf("event name = %q, want snapshot", ev.name)
	}
	var items []struct {
		ID       uuid.UUID `json:"id"`
		ItemType string    `json:"item_type"`
		Stock    *int      `json:"stock"`
	}
	err := json.Unmarshal([]byte(ev.data), &items)
	if err != nil {
		t.Fatalf("snapshot data: %v", err)
	}
	if len(items) != 3 {
		t.Fatalf("snapshot items = %d, want 3", len(items))
	}
	byID := map[uuid.UUID]struct {
		ItemType string
		Stock    *int
	}{}
	for _, it := range items {
		byID[it.ID] = struct {
			ItemType string
			Stock    *int
		}{it.ItemType, it.Stock}
	}
	if s := byID[ticketID]; s.ItemType != "ticket" || s.Stock == nil || *s.Stock != 7 {
		t.Fatalf("ticket = %+v, want ticket stock 7 (10-3)", s)
	}
	if s := byID[variantID]; s.ItemType != "product" || s.Stock != nil {
		t.Fatalf("unlimited variant = %+v, want product stock null", s)
	}
	if s := byID[occurrenceID]; s.ItemType != "program_occurrence" || s.Stock == nil || *s.Stock != 5 {
		t.Fatalf("occurrence = %+v, want program_occurrence stock 5", s)
	}
}

func TestSSEDeltaRecomputesStockFromDB(t *testing.T) {
	editionID := uuid.New()
	ops, avail, notifier := newStore(t, editionID)
	itemID := uuid.New()
	// Publisher says nothing about numbers — the relay must re-query.
	avail.set([]models.ItemAvailability{
		{ItemType: models.PurchaseItemTypeTicket, ItemID: itemID, BaseQuantity: new(int(10)), ReservedQuantity: 0},
	})

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ops.ServeStream(w, r, editionID)
	}))
	t.Cleanup(ts.Close)

	c := dialSSE(t, ts, editionID)
	if ev, ok := c.readEvent(t); !ok || ev.name != "snapshot" {
		t.Fatal("no snapshot first")
	}

	// A purchase reserves 4 units → the relay re-reads and broadcasts the
	// recomputed number (10-4=6), never a publisher-provided value.
	avail.set([]models.ItemAvailability{
		{ItemType: models.PurchaseItemTypeTicket, ItemID: itemID, BaseQuantity: new(int(10)), ReservedQuantity: 4},
	})
	raw, _ := json.Marshal(notify.StockNotification{
		Kind:      notify.KindStock,
		EditionID: editionID,
		ItemIDs:   []uuid.UUID{itemID},
	})
	notifier.fire(string(raw))

	ev, ok := c.readEvent(t)
	if !ok {
		t.Fatal("no delta event")
	}
	if ev.name != "stock" {
		t.Fatalf("event name = %q, want stock", ev.name)
	}
	var delta struct {
		ID       uuid.UUID `json:"id"`
		ItemType string    `json:"item_type"`
		Stock    int       `json:"stock"`
	}
	err := json.Unmarshal([]byte(ev.data), &delta)
	if err != nil {
		t.Fatalf("delta data: %v", err)
	}
	if delta.ID != itemID || delta.ItemType != "ticket" || delta.Stock != 6 {
		t.Fatalf("delta = %+v, want ticket stock 6 (recomputed from DB)", delta)
	}
}

// headToMap flattens raw response head lines ("Name: value") into a map,
// joining repeated headers with ", ".
func headToMap(head []string) map[string]string {
	got := map[string]string{}
	for _, line := range head {
		if name, value, ok := strings.Cut(line, ": "); ok {
			if prev, seen := got[name]; seen {
				got[name] = prev + ", " + value
			} else {
				got[name] = value
			}
		}
	}
	return got
}

// TestSSECORSHeaders pins the fix for the cross-origin stream: the harness
// CORS middleware writes Allow-Origin/Allow-Credentials/Vary into the
// ResponseWriter before the handler runs, but the SSE handler hijacks the
// connection (discarding the ResponseWriter) and emits a hand-built head.
// Those headers must survive into the raw head or the browser blocks the
// EventSource. Origin absent → no CORS headers emitted (plain clients still
// work).
func TestSSECORSHeaders(t *testing.T) {
	editionID := uuid.New()
	ops, avail, _ := newStore(t, editionID)
	avail.set(nil)

	// Simulate rs/cors handleActualRequest with AllowCredentials: true.
	origin := "http://localhost:3002"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Add("Vary", "Origin")
		ops.ServeStream(w, r, editionID)
	}))
	t.Cleanup(ts.Close)

	c, head := dialSSEHead(t, ts, editionID, "")
	defer c.conn.Close()
	got := headToMap(head)
	if v := got["Access-Control-Allow-Origin"]; v != origin {
		t.Fatalf("Access-Control-Allow-Origin = %q, want %q", v, origin)
	}
	if v := got["Access-Control-Allow-Credentials"]; v != "true" {
		t.Fatalf("Access-Control-Allow-Credentials = %q, want true", v)
	}
	if !strings.Contains(got["Vary"], "Origin") {
		t.Fatalf("Vary = %q, want it to include Origin", got["Vary"])
	}
	if v := got["Content-Type"]; !strings.HasPrefix(v, "text/event-stream") {
		t.Fatalf("Content-Type = %q, want text/event-stream", v)
	}

	// And the stream still delivers events with the CORS head in place.
	if ev, ok := c.readEvent(t); !ok || ev.name != "snapshot" {
		t.Fatalf("no snapshot event after CORS head: %+v ok=%v", ev, ok)
	}
}

// TestSSECORSViaHarnessMiddleware proves the env-configured CORS settings
// reach the hijacked stream through the real middleware the harness applies
// (fun middlewares CORS = rs/cors, AllowCredentials: true — the same
// mws.CORS config httpserver.stack builds from CORS_ALLOWED_ORIGINS). The
// handler never reads the config itself: it reproduces whatever the
// middleware stamped on the ResponseWriter before the hijack. A disallowed
// origin must get no Allow-Origin header (the browser still blocks it).
func TestSSECORSViaHarnessMiddleware(t *testing.T) {
	editionID := uuid.New()
	ops, avail, _ := newStore(t, editionID)
	avail.set(nil)

	cors := mws.CORS(mws.CORSConfig{
		AllowedOrigins:   []string{"http://localhost:3002"}, // CORS_ALLOWED_ORIGINS in .env
		AllowCredentials: true,
	})
	ts := httptest.NewServer(cors(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ops.ServeStream(w, r, editionID)
	})))
	t.Cleanup(ts.Close)

	// Allowed origin → the configured setting appears on the stream head.
	c, head := dialSSEHead(t, ts, editionID, "Origin: http://localhost:3002\r\n")
	defer c.conn.Close()
	got := headToMap(head)
	if v := got["Access-Control-Allow-Origin"]; v != "http://localhost:3002" {
		t.Fatalf("Access-Control-Allow-Origin = %q, want the configured origin", v)
	}
	if v := got["Access-Control-Allow-Credentials"]; v != "true" {
		t.Fatalf("Access-Control-Allow-Credentials = %q, want true", v)
	}
	if !strings.Contains(got["Vary"], "Origin") {
		t.Fatalf("Vary = %q, want it to include Origin", got["Vary"])
	}
	if ev, ok := c.readEvent(t); !ok || ev.name != "snapshot" {
		t.Fatalf("no snapshot after real CORS middleware: %+v ok=%v", ev, ok)
	}

	// Disallowed origin → the middleware stamps nothing, so the hijacked
	// head carries no Allow-Origin and the browser keeps blocking.
	c2, head2 := dialSSEHead(t, ts, editionID, "Origin: http://evil.example\r\n")
	defer c2.conn.Close()
	if v := headToMap(head2)["Access-Control-Allow-Origin"]; v != "" {
		t.Fatalf("disallowed origin got Access-Control-Allow-Origin = %q, want none", v)
	}
	if ev, ok := c2.readEvent(t); !ok || ev.name != "snapshot" {
		t.Fatalf("no snapshot for disallowed-origin client: %+v ok=%v", ev, ok)
	}
}

func TestSSEKeepalive(t *testing.T) {
	old := store.SseKeepaliveInterval()
	store.SetSseKeepaliveInterval(200 * time.Millisecond)
	defer store.SetSseKeepaliveInterval(old)

	editionID := uuid.New()
	ops, avail, _ := newStore(t, editionID)
	avail.set(nil)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ops.ServeStream(w, r, editionID)
	}))
	t.Cleanup(ts.Close)

	c := dialSSE(t, ts, editionID)
	if _, ok := c.readEvent(t); !ok {
		t.Fatal("no snapshot")
	}

	// The keepalive comment line arrives ~200ms later.
	_ = c.conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	line, err := c.rd.ReadString('\n')
	if err != nil {
		t.Fatalf("keepalive read: %v", err)
	}
	if !strings.HasPrefix(line, ": ping") {
		t.Fatalf("keepalive line = %q, want ': ping'", line)
	}
}

func TestSSEUnknownEdition404(t *testing.T) {
	ops, _, _ := newStore(t, uuid.New())
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ops.ServeStream(w, r, uuid.New())
	}))
	t.Cleanup(ts.Close)

	resp, err := http.NewRequestWithContext(context.Background(), http.MethodGet, ts.URL+"/editions/"+uuid.NewString()+"/store/stream", nil)
	if err != nil {
		t.Fatalf("req: %v", err)
	}
	res, err := http.DefaultClient.Do(resp)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer func() { _ = res.Body.Close() }()
	if res.StatusCode != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", res.StatusCode)
	}
}
