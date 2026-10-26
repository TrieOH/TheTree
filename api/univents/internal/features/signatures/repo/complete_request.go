package repo

import (
	"context"
	"lib/database"
	"univents/errs"
	"univents/internal/sqlc"

	"github.com/google/uuid"
)

func (repo *Repo) CompleteRequest(ctx context.Context, id, signatureID uuid.UUID) error {
	rows, err := database.Queries(ctx, repo.q).CompleteSignatureRequest(ctx, sqlc.CompleteSignatureRequestParams{
		ID:          id,
		SignatureID: &signatureID,
	})
	if err != nil {
		return repo.dbe(err)
	}
	if rows == 0 {
		return &database.NotFoundError{Named: errs.ErrSignatureRequestCompleted}
	}
	return nil
}
