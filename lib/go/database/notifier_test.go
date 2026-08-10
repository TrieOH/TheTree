package database_test

import (
	"context"
	"testing"
	"time"

	"lib/database"
	"lib/testdb"

	"github.com/jackc/pgx/v5"
)

func TestNotifierRoundTrip(t *testing.T) {
	pool := testdb.Postgres(t, "")
	dsn := pool.Config().ConnString()

	n := database.NewNotifier(dsn)
	defer n.Close()

	ctx := context.Background()
	got := make(chan string, 1)
	n.Subscribe("test_channel", func(payload string) {
		got <- payload
	})

	// The listen loop connects asynchronously; wait until it is live before
	// notifying (a notification sent before LISTEN is lost by design).
	waitForListener(t, dsn)

	err := n.Notify(ctx, "test_channel", "hello")
	if err != nil {
		t.Fatalf("Notify: %v", err)
	}

	select {
	case p := <-got:
		if p != "hello" {
			t.Fatalf("payload = %q, want hello", p)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("notification not received")
	}
}

func TestNotifierMultipleHandlers(t *testing.T) {
	pool := testdb.Postgres(t, "")
	dsn := pool.Config().ConnString()

	n := database.NewNotifier(dsn)
	defer n.Close()

	ctx := context.Background()
	a := make(chan string, 1)
	b := make(chan string, 1)
	n.Subscribe("fan_channel", func(p string) { a <- p })
	n.Subscribe("fan_channel", func(p string) { b <- p })

	waitForListener(t, dsn)

	err := n.Notify(ctx, "fan_channel", "x")
	if err != nil {
		t.Fatalf("Notify: %v", err)
	}
	for name, ch := range map[string]chan string{"a": a, "b": b} {
		select {
		case <-ch:
		case <-time.After(5 * time.Second):
			t.Fatalf("handler %s not called", name)
		}
	}
}

// waitForListener waits until the notifier's LISTEN connection is live by
// probing pg_stat_activity for a client backend issued a LISTEN statement.
func waitForListener(t *testing.T, dsn string) {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for {
		conn, err := pgxConnect(dsn)
		if err == nil {
			var listening int
			err = conn.QueryRow(context.Background(),
				`SELECT count(*) FROM pg_stat_activity
				 WHERE backend_type = 'client backend'
				   AND query ILIKE 'LISTEN %'`).
				Scan(&listening)
			_ = conn.Close(context.Background())
			if err == nil && listening > 0 {
				return
			}
		}
		if time.Now().After(deadline) {
			t.Fatal("listen connection never established")
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func pgxConnect(dsn string) (*pgx.Conn, error) {
	return pgx.Connect(context.Background(), dsn)
}

// TestNotifierReconnectsAfterListenConnDrop pins the reconnect loop: when
// the dedicated LISTEN connection dies (network blip, server restart), the
// notifier re-establishes it and keeps delivering.
func TestNotifierReconnectsAfterListenConnDrop(t *testing.T) {
	pool := testdb.Postgres(t, "")
	dsn := pool.Config().ConnString()

	n := database.NewNotifier(dsn)
	defer n.Close()

	ctx := context.Background()
	got := make(chan string, 1)
	n.Subscribe("reconnect_channel", func(payload string) {
		got <- payload
	})

	waitForListener(t, dsn)

	// Confirm round-trip works before the drop.
	err := n.Notify(ctx, "reconnect_channel", "before")
	if err != nil {
		t.Fatalf("Notify: %v", err)
	}
	select {
	case p := <-got:
		if p != "before" {
			t.Fatalf("payload = %q, want before", p)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("initial notification not received")
	}

	// Kill the listen connection; the loop reconnects + re-LISTENs.
	n.DropListenConn()

	// Notify may race the reconnect (fire-and-forget drops payloads with no
	// listener), so retry until a delivery lands.
	deadline := time.Now().Add(10 * time.Second)
	for {
		err := n.Notify(ctx, "reconnect_channel", "after")
		if err != nil {
			t.Fatalf("Notify: %v", err)
		}
		select {
		case p := <-got:
			if p != "after" {
				t.Fatalf("payload = %q, want after", p)
			}
			return // reconnected and delivered
		case <-time.After(100 * time.Millisecond):
			if time.Now().After(deadline) {
				t.Fatal("notification not received after reconnect")
			}
		}
	}
}
