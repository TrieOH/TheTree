package repos

import (
	"IdentityX/models"
	"context"
	"lib/database"
	"lib/xslices"

	"github.com/google/uuid"
)

func (repo *repo) GetActiveSigningKeys(ctx context.Context, projectID *uuid.UUID) ([]models.ActiveSigningKey, error) {
	ctx, span := database.Span(ctx, repo.tracer, "GetActiveSigningKeys")
	defer span.End()
	sqlcKeys, err := database.Queries(ctx, repo.q).GetActiveSigningKeys(ctx, projectID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(sqlcKeys, mapToActiveSigningKey), nil
}
