package repos

import (
	"context"

	"lib/database"

	"github.com/google/uuid"
)

func (repo *repo) Delete(ctx context.Context, id uuid.UUID) error {
	ctx, span := database.Span(ctx, repo.tracer, "FieldRepo.Delete")
	defer span.End()
	err := database.Queries(ctx, repo.q).DeleteField(ctx, id)
	if err != nil {
		return repo.dbe(err)
	}
	return nil
}
