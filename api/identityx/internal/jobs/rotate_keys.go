package jobs

import (
	"context"

	"IdentityX/internal/keys"

	"github.com/riverqueue/river"
)

// RotateKeysArgs drives the periodic Key-lifecycle sweep: it runs the
// keys module's EnsureAll on the ROTATE_KEYS_JOB_DURATION schedule, so
// every scope's keys are provisioned, rotated, and swept even when the
// service stays up for longer than a key lifetime. Boot also runs
// EnsureAll inline; this worker is the ongoing heartbeat (idempotent —
// Ensure is safe to run overlapping).
type RotateKeysArgs struct{}

func (RotateKeysArgs) Kind() string { return "keys.rotate" }

type RotateKeysWorker struct {
	river.WorkerDefaults[RotateKeysArgs]

	keys *keys.Manager
}

func NewRotateKeysWorker(keys *keys.Manager) *RotateKeysWorker {
	return &RotateKeysWorker{keys: keys}
}

func (w *RotateKeysWorker) Work(ctx context.Context, _ *river.Job[RotateKeysArgs]) error {
	return w.keys.EnsureAll(ctx)
}
