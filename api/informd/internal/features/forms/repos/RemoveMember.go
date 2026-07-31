package repos

import (
	"Informd/internal/sqlc"
	"context"

	"lib/database"
	"lib/telemetry"

	"github.com/google/uuid"
)

func (repo *Repo) RemoveMember(ctx context.Context, userID, formID uuid.UUID) error {
	ctx, span := telemetry.StartSpan(ctx, "RemoveMember")
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
