package repo

import (
	"Informd/internal/sqlc"
	"context"

	"Informd/models"
	"lib/database"

	"github.com/google/uuid"
)

func (repo *Repo) GetMember(ctx context.Context, userID, formID uuid.UUID) (*models.FormMember, error) {
	sqlcMember, err := database.Queries(ctx, repo.q).GetFormMember(ctx, sqlc.GetFormMemberParams{
		UserID: userID,
		FormID: formID,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapFormMember(sqlcMember)), nil
}
