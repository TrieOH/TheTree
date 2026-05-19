package forms

import (
	"Informd/internal/database/sqlc"
	"Informd/models"
	"context"
	"lib/database"

	"github.com/google/uuid"
)

func (repo *repo) GetMember(ctx context.Context, userID, formID uuid.UUID) (*models.FormMember, error) {
	ctx, span := repo.tracer.Start(ctx, "GetMember")
	defer span.End()
	sqlcMember, err := database.Queries(ctx, repo.q).GetFormMember(ctx, sqlc.GetFormMemberParams{
		UserID: userID,
		FormID: formID,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapFormMember(sqlcMember)), nil
}
