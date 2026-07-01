package repos

import (
	"IdentityX/internal/database/sqlc"
	"context"
	"lib/database"

	"github.com/google/uuid"
)

func (repo *repo) AssignToApiKey(ctx context.Context, apiKeyID uuid.UUID, capabilityIDs []uuid.UUID, assignedBy uuid.UUID) error {
	ctx, span := database.Span(ctx, repo.tracer, "AssignToApiKey")
	defer span.End()
	err := database.Queries(ctx, repo.q).AssignCapabilitiesToApiKey(ctx, sqlc.AssignCapabilitiesToApiKeyParams{
		ApiKeyID:      apiKeyID,
		CapabilityIds: capabilityIDs,
		AssignedBy:    &assignedBy,
	})
	return repo.dbe(err)
}
