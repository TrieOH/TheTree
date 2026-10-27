package repo

import (
	"IdentityX/models"
	"context"

	"github.com/google/uuid"
)

func (repo *Repo) GetRole(ctx context.Context, actorID, orgID uuid.UUID) (models.OrganizationRole, error) {
	member, err := repo.GetMember(ctx, actorID, orgID)
	if err != nil {
		return "", err
	}
	return member.Role, nil
}
