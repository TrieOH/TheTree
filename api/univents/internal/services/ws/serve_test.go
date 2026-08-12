package ws_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"univents/internal/services/notify"
	"univents/internal/services/ws"
	"univents/models"

	"github.com/MintzyG/fun"
	"github.com/coder/websocket"
	"github.com/google/uuid"
)

// ── token issuance ───────────────────────────────────────────────────────

func TestIssueTokenOwnerOnly(t *testing.T) {
	ops, tokens, purchases, _ := newOps(t)
	p := seedPending(purchases)

	raw, expires, err := ops.IssueToken(context.Background(), p.ID, p.PurchaserID)
	if err != nil {
		t.Fatalf("issue: %v", err)
	}
	if raw == "" {
		t.Fatal("empty raw token")
	}
	if d := time.Until(expires); d > 11*time.Minute || d < 9*time.Minute {
		t.Fatalf("expiry = %v, want ~10min from now", expires)
	}
	// Only the hash is stored; the raw token must not be.
	if tokens.tokens[raw] != nil {
		t.Fatal("raw token stored — only the hash should persist")
	}
}

func TestIssueTokenNonOwnerNotFound(t *testing.T) {
	ops, _, purchases, _ := newOps(t)
	p := seedPending(purchases)

	_, _, err := ops.IssueToken(context.Background(), p.ID, uuid.New())
	if err == nil || !fun.Is(err, fun.CodeNotFound) {
		t.Fatalf("non-owner issue err = %v, want NOT_FOUND (no existence leak)", err)
	}
}

func TestIssueTokenMissingPurchaseNotFound(t *testing.T) {
	ops, _, _, _ := newOps(t) //nolint:dogsled // only the operations are needed
	_, _, err := ops.IssueToken(context.Background(), uuid.New(), uuid.New())
	if err == nil || !fun.Is(err, fun.CodeNotFound) {
		t.Fatalf("missing purchase err = %v, want NOT_FOUND", err)
	}
}

// ── socket handshake + snapshot ──────────────────────────────────────────

func newServer(t *testing.T, ops *ws.Operations) *httptest.Server {
	t.Helper()
	ts := httptest.NewServer(http.HandlerFunc(ops.ServeWS))
	t.Cleanup(ts.Close)
	return ts
}

// dialWS opens a socket with the given token, returning the conn and the
// HTTP status of the handshake (101 on success; 400/401 on rejection).
func dialWS(t *testing.T, ts *httptest.Server, token string) (*websocket.Conn, int, error) {
	t.Helper()
	url := "ws" + strings.TrimPrefix(ts.URL, "http") + "/ws?token=" + token
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	conn, resp, err := websocket.Dial(ctx, url, nil)
	status := 0
	if resp != nil {
		status = resp.StatusCode
		if resp.Body != nil {
			_ = resp.Body.Close()
		}
	}
	return conn, status, err
}

// closeConn closes the socket at test teardown (errcheck-friendly).
func closeConn(t *testing.T, conn *websocket.Conn) {
	t.Helper()
	t.Cleanup(func() { _ = conn.Close(websocket.StatusNormalClosure, "") })
}

type frame struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

// readFrame reads one message with a bounded context.
func readFrame(t *testing.T, conn *websocket.Conn) (frame, error) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, data, err := conn.Read(ctx)
	if err != nil {
		return frame{}, err
	}
	var f frame
	err = json.Unmarshal(data, &f)
	if err != nil {
		t.Fatalf("bad frame json: %v (%s)", err, data)
	}
	return f, nil
}

func TestServeWSSnapshotOnConnect(t *testing.T) {
	ops, _, purchases, _ := newOps(t)
	p := seedPending(purchases)

	raw, _, err := ops.IssueToken(context.Background(), p.ID, p.PurchaserID)
	if err != nil {
		t.Fatalf("issue: %v", err)
	}
	ts := newServer(t, ops)
	conn, status, err := dialWS(t, ts, raw)
	if err != nil {
		if status != 0 {
			t.Fatalf("dial: %v (status %d)", err, status)
		}
		t.Fatalf("dial: %v", err)
	}
	closeConn(t, conn)

	f, err := readFrame(t, conn)
	if err != nil {
		t.Fatalf("first frame: %v", err)
	}
	if f.Type != "purchase.snapshot" {
		t.Fatalf("first frame type = %q, want purchase.snapshot", f.Type)
	}
	var snap struct {
		PurchaseID   uuid.UUID  `json:"purchase_id"`
		Status       string     `json:"status"`
		IntentID     *uuid.UUID `json:"intent_id"`
		IntentStatus *string    `json:"intent_status"`
		TotalCents   int64      `json:"total_cents"`
		Items        []struct {
			ItemType string    `json:"item_type"`
			ItemID   uuid.UUID `json:"item_id"`
			Quantity int       `json:"quantity"`
		} `json:"items"`
	}
	err = json.Unmarshal(f.Payload, &snap)
	if err != nil {
		t.Fatalf("snapshot payload: %v", err)
	}
	if snap.PurchaseID != p.ID || snap.Status != string(models.PurchaseStatusPending) {
		t.Fatalf("snapshot = %+v, want purchase %s pending", snap, p.ID)
	}
	if snap.TotalCents != 1234 || len(snap.Items) != 1 || snap.Items[0].Quantity != 2 {
		t.Fatalf("snapshot items = %+v, want the seeded item", snap.Items)
	}
	if snap.IntentStatus == nil || *snap.IntentStatus != "pending" {
		t.Fatalf("intent_status = %v, want pending (purchase pending + intent)", snap.IntentStatus)
	}
}

func TestServeWSSnapshotResolvedPurchase(t *testing.T) {
	// A terminal purchase (already approved): the snapshot carries the
	// resolved state and no intent fetch happens — reconnect restores
	// context from the snapshot alone.
	ops, _, purchases, _ := newOps(t)
	p := &models.Purchase{
		EditionID:   uuid.New(),
		PurchaserID: uuid.New(),
		Status:      models.PurchaseStatusApproved,
		TotalCents:  500,
		Currency:    "BRL",
		ExpiresAt:   time.Now().Add(-time.Minute),
		CreatedAt:   time.Now(),
	}
	purchases.seed(p, nil)

	raw, _, err := ops.IssueToken(context.Background(), p.ID, p.PurchaserID)
	if err != nil {
		t.Fatalf("issue: %v", err)
	}
	ts := newServer(t, ops)
	conn, _, err := dialWS(t, ts, raw)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	closeConn(t, conn)

	f, err := readFrame(t, conn)
	if err != nil {
		t.Fatalf("snapshot: %v", err)
	}
	var snap struct {
		Status string `json:"status"`
	}
	err = json.Unmarshal(f.Payload, &snap)
	if err != nil {
		t.Fatalf("snapshot payload: %v", err)
	}
	if snap.Status != string(models.PurchaseStatusApproved) {
		t.Fatalf("snapshot status = %q, want approved", snap.Status)
	}
}

func TestServeWSRejectsBadTokens(t *testing.T) {
	ops, _, purchases, _ := newOps(t)
	p := seedPending(purchases)
	ts := newServer(t, ops)

	// Missing token → 400 (no upgrade).
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, resp, err := websocket.Dial(ctx, "ws"+strings.TrimPrefix(ts.URL, "http")+"/ws", nil)
	if resp != nil && resp.Body != nil {
		_ = resp.Body.Close()
	}
	if err == nil {
		t.Fatal("missing token dial succeeded")
	}
	if resp == nil || resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("missing token status = %v, want 400", resp)
	}

	// Unknown token → 401.
	_, status, err := dialWS(t, ts, "not-a-real-token")
	if err == nil {
		t.Fatal("unknown token dial succeeded")
	}
	if status != http.StatusUnauthorized {
		t.Fatalf("unknown token status = %v, want 401", status)
	}

	// Consumed (one-time) token → 401 on reuse.
	raw, _, err := ops.IssueToken(context.Background(), p.ID, p.PurchaserID)
	if err != nil {
		t.Fatalf("issue: %v", err)
	}
	first, _, err := dialWS(t, ts, raw)
	if err != nil {
		t.Fatalf("first dial: %v", err)
	}
	_ = first.Close(websocket.StatusNormalClosure, "")
	time.Sleep(50 * time.Millisecond) // let the consume land
	_, status, err = dialWS(t, ts, raw)
	if err == nil {
		t.Fatal("reused token dial succeeded")
	}
	if status != http.StatusUnauthorized {
		t.Fatalf("reused token status = %v, want 401 (one-time)", status)
	}
}

// ── hub: fan-out + close-on-terminal ─────────────────────────────────────

func purchasePayload(purchaseID uuid.UUID, status string) string {
	raw, _ := json.Marshal(notify.PurchaseNotification{
		Kind:       notify.KindPurchase,
		EditionID:  uuid.New(),
		PurchaseID: purchaseID,
		Status:     status,
	})
	return string(raw)
}

func TestHubFanOutAndCloseOnTerminal(t *testing.T) {
	ops, _, purchases, notifier := newOps(t)
	p := seedPending(purchases)
	ts := newServer(t, ops)

	token1, _, _ := ops.IssueToken(context.Background(), p.ID, p.PurchaserID)
	token2, _, _ := ops.IssueToken(context.Background(), p.ID, p.PurchaserID)
	c1, _, err := dialWS(t, ts, token1)
	if err != nil {
		t.Fatalf("dial 1: %v", err)
	}
	closeConn(t, c1)
	c2, _, err := dialWS(t, ts, token2)
	if err != nil {
		t.Fatalf("dial 2: %v", err)
	}
	closeConn(t, c2)

	// Both clients get the snapshot first.
	_, err = readFrame(t, c1)
	if err != nil {
		t.Fatalf("c1 snapshot: %v", err)
	}
	_, err = readFrame(t, c2)
	if err != nil {
		t.Fatalf("c2 snapshot: %v", err)
	}

	notifier.fire(purchasePayload(p.ID, string(models.PurchaseStatusApproved)))

	for _, c := range []*websocket.Conn{c1, c2} {
		first, err := readFrame(t, c)
		if err != nil {
			t.Fatalf("fan-out frame: %v", err)
		}
		if first.Type != "intent.updated" {
			t.Fatalf("first live frame = %q, want intent.updated", first.Type)
		}
		var evt struct {
			PurchaseID uuid.UUID `json:"purchase_id"`
			IntentID   uuid.UUID `json:"intent_id"`
			Status     string    `json:"status"`
		}
		err = json.Unmarshal(first.Payload, &evt)
		if err != nil {
			t.Fatalf("intent.updated payload: %v", err)
		}
		if evt.PurchaseID != p.ID || evt.Status != "succeeded" {
			t.Fatalf("intent.updated = %+v, want purchase %s succeeded", evt, p.ID)
		}

		second, err := readFrame(t, c)
		if err != nil {
			t.Fatalf("terminal frame: %v", err)
		}
		if second.Type != "purchase.confirmed" {
			t.Fatalf("second frame = %q, want purchase.confirmed", second.Type)
		}

		// The terminal event closes the socket (server sends a Close frame).
		var closeErr websocket.CloseError
		err = connReadClose(c, &closeErr)
		if err != nil {
			t.Fatalf("expected close after terminal: %v", err)
		}
	}
}

// connReadClose drains frames until the server's close frame arrives.
func connReadClose(c *websocket.Conn, closeErr *websocket.CloseError) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	for {
		_, _, err := c.Read(ctx)
		if err == nil {
			continue // data frame — keep draining until the close frame
		}
		if errors.As(err, closeErr) {
			return nil
		}
		return err
	}
}

func TestHubCancelledFrameCarriesReason(t *testing.T) {
	// The cancelled frame carries the provider's status_detail via GetIntent
	// — and still goes out when Payssage is unreachable (best-effort).
	detail := "insufficient_funds"
	ops, _, purchases, notifier := newOpsWithIntents(t, &fakeIntents{
		status: "cancelled",
		detail: &detail,
	})
	p := seedPending(purchases)
	ts := newServer(t, ops)

	token, _, _ := ops.IssueToken(context.Background(), p.ID, p.PurchaserID)
	conn, _, err := dialWS(t, ts, token)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	closeConn(t, conn)
	_, err = readFrame(t, conn)
	if err != nil {
		t.Fatalf("snapshot: %v", err)
	}

	notifier.fire(purchasePayload(p.ID, string(models.PurchaseStatusCancelled)))

	first, err := readFrame(t, conn)
	if err != nil {
		t.Fatalf("intent.updated: %v", err)
	}
	if first.Type != "intent.updated" {
		t.Fatalf("frame = %q, want intent.updated", first.Type)
	}
	second, err := readFrame(t, conn)
	if err != nil {
		t.Fatalf("purchase.cancelled: %v", err)
	}
	if second.Type != "purchase.cancelled" {
		t.Fatalf("frame = %q, want purchase.cancelled", second.Type)
	}
	var payload struct {
		StatusDetail *string `json:"status_detail"`
	}
	err = json.Unmarshal(second.Payload, &payload)
	if err != nil {
		t.Fatalf("cancelled payload: %v", err)
	}
	if payload.StatusDetail == nil || *payload.StatusDetail != detail {
		t.Fatalf("status_detail = %v, want %q", payload.StatusDetail, detail)
	}
}

func TestHubCancelledWithoutIntentDetail(t *testing.T) {
	// GetIntent failing must not stall the fan-out: the cancelled frame goes
	// out without the detail.
	ops, _, purchases, notifier := newOpsWithIntents(t, &fakeIntents{err: errors.New("payssage down")})
	p := seedPending(purchases)
	ts := newServer(t, ops)

	token, _, _ := ops.IssueToken(context.Background(), p.ID, p.PurchaserID)
	conn, _, err := dialWS(t, ts, token)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	closeConn(t, conn)
	_, err = readFrame(t, conn)
	if err != nil {
		t.Fatalf("snapshot: %v", err)
	}

	notifier.fire(purchasePayload(p.ID, string(models.PurchaseStatusCancelled)))

	frames := []frame{}
	for len(frames) < 2 {
		f, err := readFrame(t, conn)
		if err != nil {
			break
		}
		frames = append(frames, f)
	}
	if len(frames) < 2 || frames[1].Type != "purchase.cancelled" {
		t.Fatalf("frames = %+v, want [intent.updated, purchase.cancelled] even with payssage down", frames)
	}
	var payload struct {
		StatusDetail *string `json:"status_detail"`
	}
	_ = json.Unmarshal(frames[1].Payload, &payload)
	if payload.StatusDetail != nil {
		t.Fatalf("status_detail = %v, want nil when GetIntent fails", payload.StatusDetail)
	}
}

// ── reconnect restores context ───────────────────────────────────────────

func TestReconnectRestoresContext(t *testing.T) {
	// The canonical flow: user opens the socket (pix pending), the payment
	// confirms while they're away, they reconnect with a fresh token — the
	// snapshot shows the resolved state without any resume polling.
	ops, _, purchases, notifier := newOps(t)
	p := seedPending(purchases)
	ts := newServer(t, ops)

	token1, _, _ := ops.IssueToken(context.Background(), p.ID, p.PurchaserID)
	c1, _, err := dialWS(t, ts, token1)
	if err != nil {
		t.Fatalf("dial 1: %v", err)
	}
	_, err = readFrame(t, c1)
	if err != nil {
		t.Fatalf("c1 snapshot: %v", err)
	}

	// Payment confirms while the user is away: the webhook receiver flips
	// the purchase in the DB first, then notifies — mirror both here.
	_, err = purchases.UpdateStatus(context.Background(), p.ID, models.PurchaseStatusApproved, nil)
	if err != nil {
		t.Fatalf("flip status: %v", err)
	}
	notifier.fire(purchasePayload(p.ID, string(models.PurchaseStatusApproved)))
	var closeErr websocket.CloseError
	err = connReadClose(c1, &closeErr)
	if err != nil {
		t.Fatalf("c1 close after confirm: %v", err)
	}
	_ = c1.CloseNow()

	// Reconnect: fresh token → snapshot shows approved.
	token2, _, err := ops.IssueToken(context.Background(), p.ID, p.PurchaserID)
	if err != nil {
		t.Fatalf("issue 2: %v", err)
	}
	c2, _, err := dialWS(t, ts, token2)
	if err != nil {
		t.Fatalf("dial 2: %v", err)
	}
	closeConn(t, c2)

	f, err := readFrame(t, c2)
	if err != nil {
		t.Fatalf("c2 snapshot: %v", err)
	}
	var snap struct {
		Status string `json:"status"`
	}
	err = json.Unmarshal(f.Payload, &snap)
	if err != nil {
		t.Fatalf("snapshot payload: %v", err)
	}
	if snap.Status != string(models.PurchaseStatusApproved) {
		t.Fatalf("reconnect snapshot status = %q, want approved", snap.Status)
	}
}

// ── dedupe: intent.updated fires on status change only ───────────────────

func TestFramesForDedupesSameStatus(t *testing.T) {
	ops, _, purchases, _ := newOps(t)
	p := seedPending(purchases)
	intentID := uuid.New()

	changed := ops.FramesForForTest(string(models.PurchaseStatusApproved), p.ID, intentID, true, true)
	if len(changed) != 2 || changed[0].Type != "intent.updated" || changed[1].Type != "purchase.confirmed" {
		t.Fatalf("changed frames = %+v, want [intent.updated, purchase.confirmed]", changed)
	}
	if !changed[1].Terminal {
		t.Fatal("confirmed frame must be terminal")
	}

	// Same status again (duplicate delivery) → no intent.updated.
	dup := ops.FramesForForTest(string(models.PurchaseStatusApproved), p.ID, intentID, true, false)
	if len(dup) != 1 || dup[0].Type != "purchase.confirmed" {
		t.Fatalf("duplicate frames = %+v, want [purchase.confirmed] only", dup)
	}

	// Expired → no intent.updated at all (the intent is untouched).
	expired := ops.FramesForForTest(string(models.PurchaseStatusExpired), p.ID, intentID, true, true)
	if len(expired) != 1 || expired[0].Type != "purchase.expired" {
		t.Fatalf("expired frames = %+v, want [purchase.expired]", expired)
	}
}
