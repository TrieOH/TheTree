package repo

import (
	"context"
	"univents/models"

	"github.com/google/uuid"
)

func (repo *Repo) GetRole(ctx context.Context, actorID, eventID uuid.UUID) (models.EventMemberRole, error) {
	member, err := repo.GetMember(ctx, eventID, actorID)
	if err != nil {
		return "", err
	}
	return member.Role, nil
}
