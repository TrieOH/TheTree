package repo

import (
	"IdentityX/models"
	"context"

	"github.com/google/uuid"
)

func (repo *Repo) GetRole(ctx context.Context, actorID, projectID uuid.UUID) (models.ProjectRole, error) {
	member, err := repo.GetMember(ctx, actorID, projectID)
	if err != nil {
		return "", err
	}
	return member.Role, nil
}
