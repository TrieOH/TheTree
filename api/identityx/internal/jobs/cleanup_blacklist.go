package jobs

import (
	"IdentityX/internal/sqlc"
	"context"

	"github.com/riverqueue/river"
)

type CleanupBlacklistArgs struct{}

func (CleanupBlacklistArgs) Kind() string { return "cleanup_blacklist" }

type CleanupBlacklistWorker struct {
	river.WorkerDefaults[CleanupBlacklistArgs]

	q *sqlc.Queries
}

func NewCleanupBlacklistWorker(q *sqlc.Queries) *CleanupBlacklistWorker {
	return &CleanupBlacklistWorker{q: q}
}

func (w *CleanupBlacklistWorker) Work(ctx context.Context, _ *river.Job[CleanupBlacklistArgs]) error {
	return w.q.DeleteExpiredBlacklistEntries(ctx)
}
