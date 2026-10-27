package repo

import (
	"context"
	"lib/database"
	"univents/errs"
	"univents/internal/sqlc"

	"github.com/google/uuid"
)

func (repo *Repo) CancelRequest(ctx context.Context, id uuid.UUID, reason *string) error {
	rows, err := database.Queries(ctx, repo.q).CancelSignatureRequest(ctx, sqlc.CancelSignatureRequestParams{
		ID:           id,
		StatusReason: reason,
	})
	if err != nil {
		return repo.dbe(err)
	}
	if rows == 0 {
		return &database.NotFoundError{Named: errs.ErrSignatureRequestHandled}
	}
	return nil
}
