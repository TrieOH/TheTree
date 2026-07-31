package repos

import (
	"context"

	"lib/database"
	"lib/telemetry"

	"github.com/google/uuid"
)

func (repo *Repo) DeleteSelectConfig(ctx context.Context, fieldID uuid.UUID) error {
	ctx, span := telemetry.StartSpan(ctx, "FieldRepo.DeleteSelectConfig")
	defer span.End()
	err := database.Queries(ctx, repo.q).DeleteFieldSelectConfig(ctx, fieldID)
	if err != nil {
		return repo.dbe(err)
	}
	return nil
}
