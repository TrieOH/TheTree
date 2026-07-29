package repos

import (
	"context"
	"lib/database"
	"lib/telemetry"
	"univents/internal/sqlc"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

func (repo *Repo) CancelRequest(ctx context.Context, id uuid.UUID, reason *string) error {
	ctx, span := telemetry.StartSpan(ctx, "SignaturesRepo.CancelRequest")
	defer span.End()
	rows, err := database.Queries(ctx, repo.q).CancelSignatureRequest(ctx, sqlc.CancelSignatureRequestParams{
		ID:           id,
		StatusReason: reason,
	})
	if err != nil {
		return repo.dbe(err)
	}
	if rows == 0 {
		return fun.ErrNotFound("signature request not found or already handled")
	}
	return nil
}
