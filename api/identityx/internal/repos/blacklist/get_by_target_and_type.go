package blacklist

import (
	"IdentityX/internal/sqlc"
	"IdentityX/models"
	"context"
	"lib/database"

	"lib/telemetry"
)

func (repo *Repo) GetByTargetAndType(ctx context.Context, target string, entryType models.BlacklistEntryType) (*models.BlacklistEntry, error) {
	ctx, span := telemetry.StartSpan(ctx, "GetByTargetAndType")
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
