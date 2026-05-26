package repos

import (
	"Informd/models"
	"context"
	"lib/database"
	"lib/xslices"

	"github.com/google/uuid"
)

func (repo *repo) ListDirectMembers(ctx context.Context, formID uuid.UUID) ([]models.FormMember, error) {
	ctx, span := repo.tracer.Start(ctx, "ListDirectFormMembers")
	defer span.End()
	sqlcMembers, err := database.Queries(ctx, repo.q).ListDirectFormMembers(ctx, formID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(sqlcMembers, mapFormMember), nil
}
