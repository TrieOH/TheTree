package repos

import (
	"IdentityX/internal/database/sqlc"
	"context"
	"lib/database"

	"github.com/google/uuid"
)

func (repo *repo) ValidateCapabilities(ctx context.Context, projectID *uuid.UUID, capabilities []uuid.UUID) (bool, error) {
	ctx, span := database.Span(ctx, repo.tracer, "ValidateCapabilities")
	defer span.End()
	valid, err := database.Queries(ctx, repo.q).ValidateCapabilities(ctx, sqlc.ValidateCapabilitiesParams{
		CapabilityCount: len(capabilities),
		ProjectID:       projectID,
		CapabilityIds:   capabilities,
	})
	return valid, repo.dbe(err)
}
