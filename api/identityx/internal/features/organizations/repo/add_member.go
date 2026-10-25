package repo

import (
	"IdentityX/internal/sqlc"
	"IdentityX/models"
	"context"
	"lib/database"
)

func (repo *Repo) AddMember(ctx context.Context, toCreate models.OrganizationMember) error {
	err := database.Queries(ctx, repo.q).AddOrganizationMember(ctx, sqlc.AddOrganizationMemberParams{
		ActorID:        toCreate.ActorID,
		OrganizationID: toCreate.OrganizationID,
		Role:           string(toCreate.Role),
	})
	if err != nil {
		return repo.dbe(err)
	}
	return nil
}
