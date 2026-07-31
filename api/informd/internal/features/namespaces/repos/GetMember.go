package repos

import (
	"Informd/internal/sqlc"
	"context"

	"Informd/models"
	"lib/database"
	"lib/telemetry"

	"github.com/google/uuid"
)

func (repo *Repo) GetMember(ctx context.Context, userID, namespaceID uuid.UUID) (*models.NamespaceMember, error) {
	ctx, span := telemetry.StartSpan(ctx, "GetMember")
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
