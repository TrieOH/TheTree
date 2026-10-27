package repo

import (
	"IdentityX/internal/sqlc"
	"context"
	"lib/database"

	"github.com/google/uuid"
)

func (repo *Repo) ValidateCapabilities(ctx context.Context, projectID *uuid.UUID, capabilities []uuid.UUID) (bool, error) {
	valid, err := database.Queries(ctx, repo.q).ValidateCapabilities(ctx, sqlc.ValidateCapabilitiesParams{
		CapabilityCount: len(capabilities),
		ProjectID:       projectID,
		CapabilityIds:   capabilities,
	})
	return valid, repo.dbe(err)
}
