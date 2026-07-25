package repos

import (
	"IdentityX/models"
	"context"
	"lib/database"
	"lib/xslices"

	"lib/telemetry"
)

func (repo *Repo) ListByAPIKeyPrefix(ctx context.Context, prefix string) ([]models.Capability, error) {
	ctx, span := telemetry.StartSpan(ctx, "ListByAPIKeyPrefix")
	defer span.End()

	capabilities, err := database.Queries(ctx, repo.q).ListCapabilitiesByApiKeyPrefix(ctx, prefix)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(capabilities, mapCapability), nil
}
