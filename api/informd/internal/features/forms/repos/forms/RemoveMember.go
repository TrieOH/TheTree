package forms

import (
	"Informd/internal/database/sqlc"
	"context"
	"lib/database"

	"github.com/google/uuid"
)

func (repo *repo) RemoveMember(ctx context.Context, userID, formID uuid.UUID) error {
	ctx, span := repo.tracer.Start(ctx, "RemoveMember")
	defer span.End()
	err := database.Queries(ctx, repo.q).RemoveFormMember(ctx, sqlc.RemoveFormMemberParams{
		UserID: userID,
		FormID: formID,
	})
	if err != nil {
		return repo.dbe(err)
	}
	return nil
}
