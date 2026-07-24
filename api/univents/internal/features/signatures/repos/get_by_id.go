package repos

import (
	"context"
	"lib/database"
	"univents/internal/database/sqlc"
	"univents/models"

	"github.com/google/uuid"
)

func (repo *repo) GetByID(ctx context.Context, id, editionID uuid.UUID) (*models.Signature, error) {
	ctx, span := database.Span(ctx, repo.tracer, "GetByID")
	defer span.End()
	sig, err := database.Queries(ctx, repo.q).GetSignatureByID(ctx, sqlc.GetSignatureByIDParams{
		ID:        id,
		EditionID: editionID,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapSignature(sig)), nil
}
