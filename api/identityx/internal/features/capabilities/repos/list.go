package repos

import (
	"IdentityX/models"
	"context"
	"lib/database"
	"lib/xslices"

	"github.com/google/uuid"
)

func (repo *repo) List(ctx context.Context, projectID uuid.UUID) ([]models.Capability, error) {
	ctx, span := database.Span(ctx, repo.tracer, "List")
	defer span.End()
	capabilities, err := database.Queries(ctx, repo.q).ListCapabilitiesByProject(ctx, &projectID)
	return xslices.MapSlice(capabilities, mapCapability), repo.dbe(err)
}
