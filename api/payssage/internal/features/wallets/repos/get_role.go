package repos

import (
	"context"
	"lib/database"
	"lib/telemetry"
	"payssage/internal/sqlc"
	"payssage/models"

	"github.com/google/uuid"
)

func (repo *Repo) GetRole(ctx context.Context, actorID, walletID uuid.UUID) (models.OrganizationRole, error) {
	ctx, span := telemetry.StartSpan(ctx, "WalletsRepo.GetRole")
	defer span.End()

	role, err := database.Queries(ctx, repo.q).GetWalletRole(ctx, sqlc.GetWalletRoleParams{
		ID:      walletID,
		ActorID: actorID,
	})
	if err != nil {
		return "", repo.dbe(err)
	}
	return models.OrganizationRole(role), nil
}
