package repos

import (
	"IdentityX/internal/sqlc"
	"context"
	"lib/database"

	"lib/telemetry"

	"github.com/google/uuid"
)

func (repo *Repo) AssignToAPIKey(ctx context.Context, apiKeyID uuid.UUID, capabilityIDs []uuid.UUID, assignedBy uuid.UUID) error {
	ctx, span := telemetry.StartSpan(ctx, "AssignToAPIKey")
	defer span.End()

	err := database.Queries(ctx, repo.q).AssignCapabilitiesToApiKey(ctx, sqlc.AssignCapabilitiesToApiKeyParams{
		ApiKeyID:      apiKeyID,
		CapabilityIds: capabilityIDs,
		AssignedBy:    &assignedBy,
	})
	return repo.dbe(err)
}
