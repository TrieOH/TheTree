package repos

import (
	"context"
	"lib/database"

	"github.com/google/uuid"
)

func (repo *repo) ResponsesCount(ctx context.Context, formID uuid.UUID) (int, error) {
	ctx, span := database.Span(ctx, repo.tracer, "FormRepo.ResponsesCount")
	defer span.End()
	count, err := database.Queries(ctx, repo.q).CountFormResponses(ctx, formID)
	if err != nil {
		return 0, repo.dbe(err)
	}
	return int(count), nil
}
