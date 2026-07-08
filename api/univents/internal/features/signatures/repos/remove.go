package repos

import (
	"context"
	"lib/database"
	"univents/internal/database/sqlc"

	"github.com/google/uuid"
)

func (repo *repo) Remove(ctx context.Context, id, editionID uuid.UUID) error {
	ctx, span := database.Span(ctx, repo.tracer, "Remove")
	defer span.End()
	err := database.Queries(ctx, repo.q).RemoveSignatureFromEdition(ctx, sqlc.RemoveSignatureFromEditionParams{
		ID:        id,
		EditionID: editionID,
	})
	return repo.dbe(err)
}
