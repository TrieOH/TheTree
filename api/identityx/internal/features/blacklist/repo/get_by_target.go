package repo

import (
	"IdentityX/models"
	"context"
	"lib/database"
)

func (repo *Repo) GetByTarget(ctx context.Context, target string) (*models.BlacklistEntry, error) {
	sqlcEntry, err := database.Queries(ctx, repo.q).GetBlacklistEntryByTarget(ctx, target)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapEntry(sqlcEntry)), nil
}
