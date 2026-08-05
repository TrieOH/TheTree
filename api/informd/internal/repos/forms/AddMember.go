package forms

import (
	"Informd/internal/sqlc"
	"context"

	"Informd/models"
	"lib/database"
	"lib/telemetry"
)

func (repo *Repo) AddMember(ctx context.Context, toCreate models.FormMember) error {
	ctx, span := telemetry.StartSpan(ctx, "AddMember")
	defer span.End()
	err := database.Queries(ctx, repo.q).AddFormMember(ctx, sqlc.AddFormMemberParams{
		UserID:  toCreate.UserID,
		FormID:  toCreate.FormID,
		Role:    string(toCreate.Role),
		AddedAt: toCreate.AddedAt,
		AddedBy: toCreate.AddedBy,
	})
	if err != nil {
		return repo.dbe(err)
	}
	return nil
}
