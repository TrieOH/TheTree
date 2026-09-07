package database

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
)

// fakeTx satisfies pgx.Tx by embedding the interface (methods never called —
// these tests only check whether a tx is *present* in the context).
type fakeTx struct{ pgx.Tx }

// fakeHandle is a minimal TxQueries implementation that records whether
// WithTx was invoked, mirroring how *sqlc.Queries behaves.
type fakeHandle struct{ withTx bool }

func (h *fakeHandle) WithTx(pgx.Tx) *fakeHandle { return &fakeHandle{withTx: true} }

func ctxWithTx() context.Context {
	return context.WithValue(context.Background(), TxKeyValue, fakeTx{})
}

// TestQueries_UsesTxFromContext pins the routing contract that WithoutTx
// builds on: a context carrying a transaction makes Queries hand the tx to
// the handle.
func TestQueries_UsesTxFromContext(t *testing.T) {
	base := &fakeHandle{}
	got := Queries(ctxWithTx(), base)
	if !got.withTx {
		t.Fatal("Queries with a tx in ctx must route to WithTx")
	}
}

// TestWithoutTx_FallsBackToPool is the regression guard for the badge/cert
// email goroutines: after WithoutTx, the context must no longer route onto
// the caller's transaction connection — otherwise concurrent queries on one
// pgx connection fail with "conn busy" / SQLSTATE 08P01.
func TestWithoutTx_FallsBackToPool(t *testing.T) {
	base := &fakeHandle{}
	clean := WithoutTx(ctxWithTx())
	got := Queries(clean, base)
	if got.withTx {
		t.Fatal("WithoutTx must strip the tx so Queries falls back to the pool handle")
	}

	// And the tx key must be gone entirely, not just typed differently.
	if v := clean.Value(TxKeyValue); v != nil {
		t.Fatalf("WithoutTx left a tx in ctx: %v", v)
	}
}

// TestWithoutTx_PreservesValuesAndCancellation ensures the helper only strips
// the transaction — telemetry, identity and cancellation semantics survive.
func TestWithoutTx_PreservesValuesAndCancellation(t *testing.T) {
	type key struct{}
	ctx, cancel := context.WithCancel(context.Background())
	ctx = context.WithValue(ctx, key{}, "kept")
	ctx = context.WithValue(ctx, TxKeyValue, fakeTx{})

	clean := WithoutTx(ctx)
	if got := clean.Value(key{}); got != "kept" {
		t.Fatalf("WithoutTx dropped unrelated context values: %v", got)
	}

	cancel()
	err := clean.Err()
	if err != context.Canceled {
		t.Fatalf("WithoutTx must preserve cancellation: got %v, want context.Canceled", err)
	}
}
