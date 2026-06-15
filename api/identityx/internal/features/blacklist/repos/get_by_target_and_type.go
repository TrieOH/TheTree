package repos

import (
	"IdentityX/internal/database/sqlc"
	"IdentityX/models"
	"context"
	"lib/database"
)

func (repo *repo) GetByTargetAndType(ctx context.Context, target string, entryType models.BlacklistEntryType) (*models.BlacklistEntry, error) {
	ctx, span := repo.tracer.Start(ctx, "Register")
	defer span.End()
	sqlcEntry, err := database.Queries(ctx, repo.q).GetBlacklistEntryByTargetAndType(ctx, sqlc.GetBlacklistEntryByTargetAndTypeParams{
		Target: target,
		Type:   string(entryType),
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapEntry(sqlcEntry)), nil
}
