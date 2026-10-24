package repo

import (
	"IdentityX/models"
	"context"
	"lib/database"
	"lib/xslices"

	"github.com/google/uuid"
)

func (repo *Repo) List(ctx context.Context, projectID uuid.UUID) ([]models.Capability, error) {
	capabilities, err := database.Queries(ctx, repo.q).ListCapabilitiesByProject(ctx, &projectID)
	return xslices.MapSlice(capabilities, mapCapability), repo.dbe(err)
}
