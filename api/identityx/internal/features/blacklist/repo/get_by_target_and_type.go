package repo

import (
	"IdentityX/internal/sqlc"
	"IdentityX/models"
	"context"
	"lib/database"
)

func (repo *Repo) GetByTargetAndType(ctx context.Context, target string, entryType models.BlacklistEntryType) (*models.BlacklistEntry, error) {
	sqlcEntry, err := database.Queries(ctx, repo.q).GetBlacklistEntryByTargetAndType(ctx, sqlc.GetBlacklistEntryByTargetAndTypeParams{
		Target: target,
		Type:   string(entryType),
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapEntry(sqlcEntry)), nil
}
