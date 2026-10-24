package repo

import (
	"context"

	"Informd/models"
	"lib/database"

	"github.com/google/uuid"
)

func (repo *Repo) ReDraft(ctx context.Context, formID uuid.UUID) (*models.Form, error) {
	sqlcForm, err := database.Queries(ctx, repo.q).DraftForm(ctx, formID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapForm(sqlcForm)), nil
}
