package repos

import (
	"context"

	"Informd/models"
	"lib/database"
	"lib/telemetry"
	"lib/xslices"

	"github.com/google/uuid"
)

func (repo *Repo) ListFromNamespaceArchived(ctx context.Context, namespaceID uuid.UUID) ([]models.Form, error) {
	ctx, span := telemetry.StartSpan(ctx, "ListFromNamespaceArchived")
	defer span.End()
	sqlcForms, err := database.Queries(ctx, repo.q).ListNamespaceArchivedForms(ctx, &namespaceID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(sqlcForms, mapForm), nil
}
