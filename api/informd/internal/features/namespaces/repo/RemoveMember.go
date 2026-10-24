package repo

import (
	"Informd/internal/sqlc"
	"context"

	"lib/database"

	"github.com/google/uuid"
)

func (repo *Repo) RemoveMember(ctx context.Context, userID, namespaceID uuid.UUID) error {
	err := database.Queries(ctx, repo.q).RemoveNamespaceMember(ctx, sqlc.RemoveNamespaceMemberParams{
		UserID:      userID,
		NamespaceID: namespaceID,
	})
	if err != nil {
		return repo.dbe(err)
	}
	return nil
}
