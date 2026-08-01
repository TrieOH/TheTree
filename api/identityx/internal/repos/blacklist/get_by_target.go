package blacklist

import (
	"IdentityX/models"
	"context"
	"lib/database"

	"lib/telemetry"
)

func (repo *Repo) GetByTarget(ctx context.Context, target string) (*models.BlacklistEntry, error) {
	ctx, span := telemetry.StartSpan(ctx, "GetByTarget")
	defer span.End()

	sqlcEntry, err := database.Queries(ctx, repo.q).GetBlacklistEntryByTarget(ctx, target)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapEntry(sqlcEntry)), nil
}
