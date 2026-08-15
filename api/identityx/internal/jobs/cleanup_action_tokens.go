package jobs

import (
	"context"

	"IdentityX/internal/tokens"

	"github.com/riverqueue/river"
)

// CleanupActionTokensArgs sweeps expired single-use action tokens so the
// anti-replay table stays bounded. Runs on the same periodic schedule as
// the blacklist cleanup.
type CleanupActionTokensArgs struct{}

func (CleanupActionTokensArgs) Kind() string { return "cleanup_action_tokens" }

type CleanupActionTokensWorker struct {
	river.WorkerDefaults[CleanupActionTokensArgs]

	actionTokens *tokens.ActionTokenManager
}

func NewCleanupActionTokensWorker(actionTokens *tokens.ActionTokenManager) *CleanupActionTokensWorker {
	return &CleanupActionTokensWorker{actionTokens: actionTokens}
}

func (w *CleanupActionTokensWorker) Work(ctx context.Context, _ *river.Job[CleanupActionTokensArgs]) error {
	return w.actionTokens.DeleteExpired(ctx)
}
