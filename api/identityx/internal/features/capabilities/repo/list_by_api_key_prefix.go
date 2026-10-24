package repo

import (
	"IdentityX/models"
	"context"
	"lib/database"
	"lib/xslices"
)

func (repo *Repo) ListByAPIKeyPrefix(ctx context.Context, prefix string) ([]models.Capability, error) {
	capabilities, err := database.Queries(ctx, repo.q).ListCapabilitiesByApiKeyPrefix(ctx, prefix)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(capabilities, mapCapability), nil
}
