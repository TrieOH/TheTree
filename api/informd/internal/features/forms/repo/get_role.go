package repo

import (
	"Informd/models"
	"context"

	"github.com/google/uuid"
)

func (repo *Repo) GetRole(ctx context.Context, actorID, formID uuid.UUID) (models.FormMemberRole, error) {
	member, err := repo.GetMember(ctx, actorID, formID)
	if err != nil {
		return "", err
	}
	return member.Role, nil
}
