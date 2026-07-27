package repos

import (
	"IdentityX/models"
	"context"
	"lib/database"
	"lib/xslices"

	"github.com/google/uuid"
	"lib/telemetry"
)

func (repo *Repo) List(ctx context.Context, projectID uuid.UUID) ([]models.Capability, error) {
	ctx, span := telemetry.StartSpan(ctx, "List")
	defer span.End()

	capabilities, err := database.Queries(ctx, repo.q).ListCapabilitiesByProject(ctx, &projectID)
	return xslices.MapSlice(capabilities, mapCapability), repo.dbe(err)
}
