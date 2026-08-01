package repos

import (
	"Informd/models"
	"context"
	"lib/telemetry"

	"github.com/google/uuid"
)

func (repo *Repo) GetRole(ctx context.Context, actorID, formID uuid.UUID) (models.FormMemberRole, error) {
	ctx, span := telemetry.StartSpan(ctx, "FormsRepo.GetRole")
	defer span.End()

	member, err := repo.GetMember(ctx, actorID, formID)
	if err != nil {
		return "", err
	}
	return member.Role, nil
}
