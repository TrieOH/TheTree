package repos

import (
	"Informd/internal/sqlc"
	"context"

	"Informd/models"
	"lib/database"
	"lib/telemetry"
)

func (repo *Repo) AddMember(ctx context.Context, toCreate models.NamespaceMember) error {
	ctx, span := telemetry.StartSpan(ctx, "AddMember")
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
