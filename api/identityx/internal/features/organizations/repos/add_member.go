package repos

import (
	"IdentityX/internal/database/sqlc"
	"IdentityX/models"
	"context"
	"lib/database"
)

func (repo *repo) AddMember(ctx context.Context, toCreate models.OrganizationMember) error {
	ctx, span := repo.tracer.Start(ctx, "AddMember")
	defer span.End()
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
