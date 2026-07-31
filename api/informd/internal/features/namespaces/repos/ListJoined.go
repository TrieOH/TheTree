package repos

import (
	"context"

	"Informd/models"
	"lib/database"
	"lib/telemetry"
	"lib/xslices"

	"github.com/google/uuid"
)

func (repo *Repo) ListJoined(ctx context.Context, userID uuid.UUID) ([]models.Namespace, error) {
	ctx, span := telemetry.StartSpan(ctx, "ListJoined")
	defer span.End()
	sqlcNamespaces, err := database.Queries(ctx, repo.q).ListJoinedNamespaces(ctx, userID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(sqlcNamespaces, mapNamespace), nil
}
