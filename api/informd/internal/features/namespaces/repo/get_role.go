package repo

import (
	"Informd/models"
	"context"

	"github.com/google/uuid"
)

func (repo *Repo) GetRole(ctx context.Context, actorID, namespaceID uuid.UUID) (models.NamespaceMemberRole, error) {
	member, err := repo.GetMember(ctx, actorID, namespaceID)
	if err != nil {
		return "", err
	}
	return member.Role, nil
}
