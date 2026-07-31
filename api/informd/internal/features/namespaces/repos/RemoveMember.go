package repos

import (
	"Informd/internal/sqlc"
	"context"

	"lib/database"
	"lib/telemetry"

	"github.com/google/uuid"
)

func (repo *Repo) RemoveMember(ctx context.Context, userID, namespaceID uuid.UUID) error {
	ctx, span := telemetry.StartSpan(ctx, "RemoveMember")
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
