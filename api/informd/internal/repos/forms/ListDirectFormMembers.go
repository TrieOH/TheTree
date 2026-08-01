package forms

import (
	"context"

	"Informd/models"
	"lib/database"
	"lib/telemetry"
	"lib/xslices"

	"github.com/google/uuid"
)

func (repo *Repo) ListDirectMembers(ctx context.Context, formID uuid.UUID) ([]models.FormMember, error) {
	ctx, span := telemetry.StartSpan(ctx, "ListDirectFormMembers")
	defer span.End()
	sqlcMembers, err := database.Queries(ctx, repo.q).ListDirectFormMembers(ctx, formID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(sqlcMembers, mapFormMember), nil
}
