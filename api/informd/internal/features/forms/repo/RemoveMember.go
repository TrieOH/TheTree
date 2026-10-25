package repo

import (
	"Informd/internal/sqlc"
	"context"

	"lib/database"

	"github.com/google/uuid"
)

func (repo *Repo) RemoveMember(ctx context.Context, userID, formID uuid.UUID) error {
	err := database.Queries(ctx, repo.q).RemoveFormMember(ctx, sqlc.RemoveFormMemberParams{
		UserID: userID,
		FormID: formID,
	})
	if err != nil {
		return repo.dbe(err)
	}
	return nil
}
