package repos

import (
	"Informd/models"
	"context"
	"lib/database"
	"lib/xslices"

	"github.com/google/uuid"
)

func (repo *repo) ListMembers(ctx context.Context, namespaceID uuid.UUID) ([]models.NamespaceMember, error) {
	ctx, span := repo.tracer.Start(ctx, "BulkGet")
	defer span.End()
	sqlcMembers, err := database.Queries(ctx, repo.q).ListNamespaceMembers(ctx, namespaceID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(sqlcMembers, mapNamespaceMember), nil
}
