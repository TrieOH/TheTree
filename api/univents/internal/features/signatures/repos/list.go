package repos

import (
	"context"
	"lib/database"
	"lib/xslices"
	"univents/models"

	"github.com/google/uuid"
)

func (repo *repo) List(ctx context.Context, editionID uuid.UUID) ([]models.Signature, error) {
	ctx, span := database.Span(ctx, repo.tracer, "List")
	defer span.End()
	sigs, err := database.Queries(ctx, repo.q).ListSignaturesFromEdition(ctx, editionID)
	return xslices.MapSlice(sigs, mapSignature), repo.dbe(err)
}
