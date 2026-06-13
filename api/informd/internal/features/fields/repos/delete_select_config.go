package repos

import (
	"context"

	"lib/database"

	"github.com/google/uuid"
)

func (repo *repo) DeleteSelectConfig(ctx context.Context, fieldID uuid.UUID) error {
	ctx, span := database.Span(ctx, repo.tracer, "FieldRepo.DeleteSelectConfig")
	defer span.End()
	err := database.Queries(ctx, repo.q).DeleteFieldSelectConfig(ctx, fieldID)
	if err != nil {
		return repo.dbe(err)
	}
	return nil
}
