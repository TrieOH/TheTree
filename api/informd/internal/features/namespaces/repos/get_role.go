package repos

import (
	"Informd/models"
	"context"
	"lib/telemetry"

	"github.com/google/uuid"
)

func (repo *Repo) GetRole(ctx context.Context, actorID, namespaceID uuid.UUID) (models.NamespaceMemberRole, error) {
	ctx, span := telemetry.StartSpan(ctx, "NamespacesRepo.GetRole")
	defer span.End()

	member, err := repo.GetMember(ctx, actorID, namespaceID)
	if err != nil {
		return "", err
	}
	return member.Role, nil
}
