package repos

import (
	"Informd/internal/database/sqlc"
	"Informd/models"
	"context"
	"lib/database"

	"github.com/google/uuid"
)

func (repo *repo) GetMember(ctx context.Context, userID, namespaceID uuid.UUID) (*models.NamespaceMember, error) {
	ctx, span := repo.tracer.Start(ctx, "GetMember")
	defer span.End()
	sqlcMember, err := database.Queries(ctx, repo.q).GetNamespaceMember(ctx, sqlc.GetNamespaceMemberParams{
		UserID:      userID,
		NamespaceID: namespaceID,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapNamespaceMember(sqlcMember)), nil
}
