package repo

import (
	"context"
	"lib/database"
	"payssage/internal/sqlc"
	"payssage/models"
)

func (repo *Repo) AddMember(ctx context.Context, toAdd models.OrganizationMember) error {
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
