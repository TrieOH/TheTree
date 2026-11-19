package ws_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"univents/models"
)

// TestHubConcurrentNotifyAndDisconnect hammers the hub with notifications
// while clients come and go on the same purchase. Guards the fan-out's map
// access: the notify path must iterate a SNAPSHOT of the client set, never
// the live map, or unregister (map write) during iteration panics / races.
// Run with -race: a data race here fails the test.
func TestHubConcurrentNotifyAndDisconnect(t *testing.T) {
	ops, _, purchases, notifier := newOps(t)
	p := seedPending(purchases)
	ts := newServer(t, ops)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Clients that connect, read a couple frames, and leave — while
	// notifications keep landing. The disconnect (unregister) must never
	// overlap the fan-out's iteration of the client set.
	var wg sync.WaitGroup
	for range 8 {
		wg.Go(func() {
			for range 6 {
				token, _, err := ops.IssueToken(ctx, p.ID, p.PurchaserID)
				if err != nil {
					return
				}
				conn, _, err := dialWS(t, ts, token)
				if err != nil {
					continue
				}
				// Read until the server closes us or the conn dies.
				for {
					_, _, err := conn.Read(ctx)
					if err != nil {
						break
					}
				}
				_ = conn.CloseNow()
			}
		})
	}

	// Notifications: each lands in a hub goroutine; between them the
	// clients above are disconnecting (unregistering) concurrently.
	for range 60 {
		notifier.fire(purchasePayload(p.ID, string(models.PurchaseStatusApproved)))
		time.Sleep(5 * time.Millisecond)
	}
	wg.Wait()
	_ = cancel
}

// TestHubSlowClientDroppedWithoutPanic forces the enqueue drop path (full
// send queue) while the notify goroutine is iterating — the exact case that
// used to panic "concurrent map iteration and map write".
func TestHubSlowClientDroppedWithoutPanic(t *testing.T) {
	ops, _, purchases, notifier := newOps(t)
	p := seedPending(purchases)
	ts := newServer(t, ops)

	token, _, err := ops.IssueToken(context.Background(), p.ID, p.PurchaserID)
	if err != nil {
		t.Fatalf("issue: %v", err)
	}
	conn, _, err := dialWS(t, ts, token)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer func() { _ = conn.CloseNow() }()

	// Drain the snapshot, then stop reading so the send queue fills and the
	// next notification's enqueue hits the drop path (unregister mid-iterate).
	_, err = readFrame(t, conn)
	if err != nil {
		t.Fatalf("snapshot: %v", err)
	}

	done := make(chan struct{})
	go func() {
		defer close(done)
		// 60 rapid notifications: queue fills (buffer 32), drops fire,
		// clients unregister — must not panic.
		for range 60 {
			notifier.fire(purchasePayload(p.ID, string(models.PurchaseStatusApproved)))
		}
	}()

	select {
	case <-done:
	case <-time.After(10 * time.Second):
		t.Fatal("notifications stalled")
	}

	// The socket should have been closed by the drop (or still be draining);
	// either way the process must be alive — a panic fails the test.
	time.Sleep(100 * time.Millisecond)
}
