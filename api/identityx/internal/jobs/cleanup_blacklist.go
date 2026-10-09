package jobs

import (
	"context"

	"IdentityX/internal/sqlc"
)

type CleanupBlacklistArgs struct{}

func (CleanupBlacklistArgs) Kind() string { return "cleanup_blacklist" }

type CleanupBlacklistWorker struct {
	q *sqlc.Queries
}

func NewCleanupBlacklistWorker(q *sqlc.Queries) *CleanupBlacklistWorker {
	return &CleanupBlacklistWorker{q: q}
}

func (w *CleanupBlacklistWorker) Work(ctx context.Context, _ CleanupBlacklistArgs) error {
	return w.q.DeleteExpiredBlacklistEntries(ctx)
}
