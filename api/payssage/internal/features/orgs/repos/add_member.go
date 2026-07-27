package repos

import (
	"context"
	"lib/database"
	"payssage/internal/sqlc"
	"payssage/models"
)

func (repo *repo) AddMember(ctx context.Context, toAdd models.OrganizationMember) error {
	ctx, span := repo.tracer.Start(ctx, "AddMember")
	defer span.End()
	err := database.Queries(ctx, repo.q).AddOrganizationMember(ctx, sqlc.AddOrganizationMemberParams{
		MemberID:       toAdd.MemberID,
		OrganizationID: toAdd.OrganizationID,
		Role:           string(toAdd.Role),
	})
	if err != nil {
		return repo.dbe(err)
	}
	return nil
}
