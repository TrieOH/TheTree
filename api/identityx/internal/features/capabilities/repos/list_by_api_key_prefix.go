package repos

import (
	"IdentityX/models"
	"context"
	"lib/database"
	"lib/xslices"
)

func (repo *repo) ListByApiKeyPrefix(ctx context.Context, prefix string) ([]models.Capability, error) {
	ctx, span := database.Span(ctx, repo.tracer, "ListByApiKeyPrefix")
	defer span.End()
	capabilities, err := database.Queries(ctx, repo.q).ListCapabilitiesByApiKeyPrefix(ctx, prefix)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(capabilities, mapCapability), nil
}
