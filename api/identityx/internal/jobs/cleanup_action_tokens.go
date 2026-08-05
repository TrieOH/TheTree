package jobs

import (
	"IdentityX/internal/sqlc"
	"context"

	"github.com/riverqueue/river"
)

// CleanupActionTokensArgs sweeps expired single-use action tokens so the
// anti-replay table stays bounded. Runs on the same periodic schedule as
// the blacklist cleanup.
type CleanupActionTokensArgs struct{}

func (CleanupActionTokensArgs) Kind() string { return "cleanup_action_tokens" }

type CleanupActionTokensWorker struct {
	river.WorkerDefaults[CleanupActionTokensArgs]

	q *sqlc.Queries
}

func NewCleanupActionTokensWorker(q *sqlc.Queries) *CleanupActionTokensWorker {
	return &CleanupActionTokensWorker{q: q}
}

func (w *CleanupActionTokensWorker) Work(ctx context.Context, _ *river.Job[CleanupActionTokensArgs]) error {
	return w.q.DeleteExpiredActionTokens(ctx)
}
