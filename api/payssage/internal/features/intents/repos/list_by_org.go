package repos

import (
	"context"
	"lib/database"
	"lib/xslices"
	"payssage/models"

	"github.com/google/uuid"
)

func (repo *Repo) ListByOrg(ctx context.Context, orgID uuid.UUID) ([]models.Intent, error) {
	ctx, span := repo.tracer.Start(ctx, "IntentRepo.ListByOrg")
	defer span.End()

	sqlcIntents, err := database.Queries(ctx, repo.q).ListIntentsByOrg(ctx, &orgID)
	if err != nil {
		return nil, repo.dbe(err)
	}

	return xslices.MapSlice(sqlcIntents, mapIntent), nil
}
