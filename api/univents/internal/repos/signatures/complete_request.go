package signatures

import (
	"context"
	"lib/database"
	"lib/telemetry"
	"univents/internal/sqlc"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

func (repo *Repo) CompleteRequest(ctx context.Context, id, signatureID uuid.UUID) error {
	ctx, span := telemetry.StartSpan(ctx, "SignaturesRepo.CompleteRequest")
	defer span.End()
	rows, err := database.Queries(ctx, repo.q).CompleteSignatureRequest(ctx, sqlc.CompleteSignatureRequestParams{
		ID:          id,
		SignatureID: &signatureID,
	})
	if err != nil {
		return repo.dbe(err)
	}
	if rows == 0 {
		return fun.ErrNotFound("signature request not found or already completed")
	}
	return nil
}
