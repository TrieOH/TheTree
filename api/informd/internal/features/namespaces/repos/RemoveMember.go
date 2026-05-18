package repos

import (
	"Informd/internal/database/sqlc"
	"context"
	"lib/database"

	"github.com/google/uuid"
)

func (repo *repo) RemoveMember(ctx context.Context, userID, namespaceID uuid.UUID) error {
	ctx, span := repo.tracer.Start(ctx, "RemoveMember")
	defer span.End()
	err := database.Queries(ctx, repo.q).RemoveNamespaceMember(ctx, sqlc.RemoveNamespaceMemberParams{
		UserID:      userID,
		NamespaceID: namespaceID,
	})
	if err != nil {
		return repo.dbe(err)
	}
	return nil
}
