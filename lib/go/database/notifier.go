package database

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
)

// Notifier is a LISTEN/NOTIFY bridge on a dedicated PostgreSQL connection.
// The connection pool cannot LISTEN directly (notifications arrive on the
// connection that issued LISTEN and would be lost in a pool), so the
// notifier owns one long-lived pgx.Conn for subscribing and reconnects it
// with backoff when it drops.
//
// Publishers (checkout, expiry, webhook receiver) call Notify — cheap,
// fire-and-forget, no listener required. Subscribers (SSE relay, WS hub,
// splits 6) call Subscribe with a handler; multiple handlers may register
// for the same channel, each is invoked with the notification payload.
type Notifier struct {
	dsn string

	mu        sync.Mutex
	listeners map[string][]func(payload string)
	conn      *pgx.Conn // dedicated listen connection; nil until started
	started   bool
	closeCh   chan struct{}
	wg        sync.WaitGroup
}

// NewNotifier builds a notifier from a Postgres DSN. No connection is made
// until the first Subscribe.
func NewNotifier(dsn string) *Notifier {
	return &Notifier{
		dsn:       dsn,
		listeners: make(map[string][]func(payload string)),
		closeCh:   make(chan struct{}),
	}
}

// Notify publishes a payload on a channel (SELECT pg_notify). Fire-and-
// forget: the payload is dropped if no connection is listening right now —
// exactly the semantics the store needs (stock deltas re-read from the DB,
// so a missed notification is just a stale snapshot, never data loss).
func (n *Notifier) Notify(ctx context.Context, channel, payload string) error {
	conn, err := pgx.Connect(ctx, n.dsn)
	if err != nil {
		return err
	}
	defer func() { _ = conn.Close(ctx) }()

	_, err = conn.Exec(ctx, "SELECT pg_notify($1, $2)", channel, payload)
	return err
}

// Subscribe registers a handler for a channel and ensures the listen
// connection is running (started lazily on first subscribe). The handler is
// invoked synchronously on the listen goroutine — keep it cheap.
func (n *Notifier) Subscribe(channel string, handler func(payload string)) {
	n.mu.Lock()
	n.listeners[channel] = append(n.listeners[channel], handler)
	if !n.started {
		n.started = true
		n.wg.Add(1)
		go n.run()
	}
	n.mu.Unlock()
}

// Close stops the listen loop and closes the dedicated connection. Idempotent.
func (n *Notifier) Close() {
	close(n.closeCh)
	n.mu.Lock()
	conn := n.conn
	n.conn = nil
	n.mu.Unlock()
	if conn != nil {
		_ = conn.Close(context.Background())
	}
	n.wg.Wait()
}

// DropListenConn force-closes the dedicated LISTEN connection so the
// reconnect loop re-establishes it. Test-only: simulates a network blip or
// server restart. Safe to call when the connection is already gone.
func (n *Notifier) DropListenConn() {
	n.mu.Lock()
	conn := n.conn
	n.conn = nil
	n.mu.Unlock()
	if conn != nil {
		_ = conn.Close(context.Background())
	}
}

// run owns the listen loop: connect, LISTEN every subscribed channel, then
// wait for notifications; on any error, reconnect with backoff. The loop
// exits when Close closes closeCh.
func (n *Notifier) run() {
	defer n.wg.Done()

	backoff := 100 * time.Millisecond
	for {
		select {
		case <-n.closeCh:
			return
		default:
		}

		conn, err := n.connectAndListen()
		if err != nil {
			slog.Warn("notifier: listen connect failed", "err", err)
			if !sleepCtx(n.closeCh, backoff) {
				return
			}
			backoff = min(backoff*2, 5*time.Second)
			continue
		}
		backoff = 100 * time.Millisecond // reset after a successful connect

		n.mu.Lock()
		n.conn = conn
		n.mu.Unlock()

		err = n.waitForNotifications(conn)
		_ = conn.Close(context.Background())

		n.mu.Lock()
		if n.conn == conn {
			n.conn = nil
		}
		n.mu.Unlock()

		// waitForNotifications only returns on a connection error, so the
		// dropped-connection warn fires whenever the conn was not closed
		// intentionally (SA4023: the error is never nil).
		if !isClosed(conn) {
			slog.Warn("notifier: listen loop dropped", "err", err)
		}
		if !sleepCtx(n.closeCh, backoff) {
			return
		}
	}
}

// connectAndListen opens the dedicated connection and LISTENs every channel
// that currently has subscribers.
func (n *Notifier) connectAndListen() (*pgx.Conn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := pgx.Connect(ctx, n.dsn)
	if err != nil {
		return nil, err
	}

	n.mu.Lock()
	channels := make([]string, 0, len(n.listeners))
	for ch := range n.listeners {
		channels = append(channels, ch)
	}
	n.mu.Unlock()

	for _, ch := range channels {
		_, err := conn.Exec(ctx, "LISTEN "+ch)
		if err != nil {
			_ = conn.Close(context.Background())
			return nil, err
		}
	}
	return conn, nil
}

// waitForNotifications blocks until the connection drops or closes, fanning
// each notification out to the channel's handlers.
func (n *Notifier) waitForNotifications(conn *pgx.Conn) error {
	for {
		notification, err := conn.WaitForNotification(context.Background())
		if err != nil {
			return err
		}
		n.mu.Lock()
		handlers := append([]func(string){}, n.listeners[notification.Channel]...)
		n.mu.Unlock()
		for _, h := range handlers {
			h(notification.Payload)
		}
	}
}

// isClosed reports whether the error came from closing the connection
// ourselves (Close), in which case the loop should exit silently.
func isClosed(conn *pgx.Conn) bool {
	return conn.IsClosed()
}

func sleepCtx(done <-chan struct{}, d time.Duration) bool {
	select {
	case <-done:
		return false
	case <-time.After(d):
		return true
	}
}
