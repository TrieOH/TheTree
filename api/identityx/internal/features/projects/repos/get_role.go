package repos

import (
	"IdentityX/models"
	"context"
	"lib/telemetry"

	"github.com/google/uuid"
)

func (repo *Repo) GetRole(ctx context.Context, actorID, projectID uuid.UUID) (models.ProjectRole, error) {
	ctx, span := telemetry.StartSpan(ctx, "ProjectsRepo.GetRole")
	defer span.End()

	member, err := repo.GetMember(ctx, actorID, projectID)
	if err != nil {
		return "", err
	}
	return member.Role, nil
}
