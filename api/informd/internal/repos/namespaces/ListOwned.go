package namespaces

import (
	"context"

	"Informd/models"
	"lib/database"
	"lib/telemetry"
	"lib/xslices"

	"github.com/google/uuid"
)

func (repo *Repo) ListOwned(ctx context.Context, userID uuid.UUID) ([]models.Namespace, error) {
	ctx, span := telemetry.StartSpan(ctx, "ListOwned")
	defer span.End()
	sqlcNamespaces, err := database.Queries(ctx, repo.q).ListOwnedNamespaces(ctx, userID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(sqlcNamespaces, mapNamespace), nil
}
