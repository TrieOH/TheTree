package repos

import (
	"context"
	"lib/database"
	"lib/telemetry"

	"github.com/google/uuid"
)

func (repo *Repo) Delete(ctx context.Context, id uuid.UUID) error {
	ctx, span := telemetry.StartSpan(ctx, "BadgeTemplatesRepo.DeleteTemplate")
	defer span.End()
	err := database.Queries(ctx, repo.q).DeleteBadgeTemplate(ctx, id)
	if err != nil {
		return repo.dbe(err)
	}

	return nil
}
