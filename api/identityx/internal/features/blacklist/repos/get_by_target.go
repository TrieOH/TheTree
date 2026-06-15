package repos

import (
	"IdentityX/models"
	"context"
	"lib/database"
)

func (repo *repo) GetByTarget(ctx context.Context, target string) (*models.BlacklistEntry, error) {
	ctx, span := repo.tracer.Start(ctx, "Register")
	defer span.End()
	sqlcEntry, err := database.Queries(ctx, repo.q).GetBlacklistEntryByTarget(ctx, target)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapEntry(sqlcEntry)), nil
}
