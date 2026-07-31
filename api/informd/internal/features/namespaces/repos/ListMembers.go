package repos

import (
	"context"

	"Informd/models"
	"lib/database"
	"lib/telemetry"
	"lib/xslices"

	"github.com/google/uuid"
)

func (repo *Repo) ListMembers(ctx context.Context, namespaceID uuid.UUID) ([]models.NamespaceMember, error) {
	ctx, span := telemetry.StartSpan(ctx, "ListMembers")
	defer span.End()
	sqlcMembers, err := database.Queries(ctx, repo.q).ListNamespaceMembers(ctx, namespaceID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(sqlcMembers, mapNamespaceMember), nil
}
