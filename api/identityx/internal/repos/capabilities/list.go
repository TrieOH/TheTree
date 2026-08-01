package capabilities

import (
	"IdentityX/models"
	"context"
	"lib/database"
	"lib/xslices"

	"lib/telemetry"

	"github.com/google/uuid"
)

func (repo *Repo) List(ctx context.Context, projectID uuid.UUID) ([]models.Capability, error) {
	ctx, span := telemetry.StartSpan(ctx, "List")
	defer span.End()

	capabilities, err := database.Queries(ctx, repo.q).ListCapabilitiesByProject(ctx, &projectID)
	return xslices.MapSlice(capabilities, mapCapability), repo.dbe(err)
}
