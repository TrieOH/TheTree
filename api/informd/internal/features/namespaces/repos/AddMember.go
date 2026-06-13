package repos

import (
	"context"

	"Informd/internal/database/sqlc"
	"Informd/models"
	"lib/database"
)

func (repo *repo) AddMember(ctx context.Context, toCreate models.NamespaceMember) error {
	ctx, span := repo.tracer.Start(ctx, "AddMember")
	defer span.End()
	err := database.Queries(ctx, repo.q).AddNamespaceMember(ctx, sqlc.AddNamespaceMemberParams{
		UserID:      toCreate.UserID,
		NamespaceID: toCreate.NamespaceID,
		Role:        string(toCreate.Role),
		AddedAt:     toCreate.AddedAt,
		AddedBy:     toCreate.AddedBy,
	})
	if err != nil {
		return repo.dbe(err)
	}
	return nil
}
