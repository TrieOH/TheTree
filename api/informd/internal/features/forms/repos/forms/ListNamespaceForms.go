package forms

import (
	"Informd/models"
	"context"
	"lib/database"
	"lib/xslices"

	"github.com/google/uuid"
)

func (repo *repo) ListFromNamespace(ctx context.Context, namespaceID uuid.UUID) ([]models.Form, error) {
	ctx, span := repo.tracer.Start(ctx, "ListMine")
	defer span.End()
	sqlcForms, err := database.Queries(ctx, repo.q).ListNamespaceForms(ctx, &namespaceID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(sqlcForms, mapForm), nil
}
