package repos

import (
	"IdentityX/models"
	"context"
	"lib/telemetry"

	"github.com/google/uuid"
)

func (repo *Repo) GetRole(ctx context.Context, actorID, orgID uuid.UUID) (models.OrganizationRole, error) {
	ctx, span := telemetry.StartSpan(ctx, "OrganizationsRepo.GetRole")
	defer span.End()

	member, err := repo.GetMember(ctx, actorID, orgID)
	if err != nil {
		return "", err
	}
	return member.Role, nil
}
