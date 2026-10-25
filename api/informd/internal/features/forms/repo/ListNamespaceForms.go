package repo

import (
	"context"

	"Informd/models"
	"lib/database"
	"lib/xslices"

	"github.com/google/uuid"
)

func (repo *Repo) ListFromNamespace(ctx context.Context, namespaceID uuid.UUID) ([]models.Form, error) {
	sqlcForms, err := database.Queries(ctx, repo.q).ListNamespaceForms(ctx, &namespaceID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(sqlcForms, mapForm), nil
}
