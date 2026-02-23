package queries

import (
	"context"
	"univents/internal/core/domain"

	"github.com/google/uuid"
)

func (uc *QueryService) ListEditions(ctx context.Context, eventID uuid.UUID) (out []domain.Edition, err error) { // FIXME Pagination
	ctx, span := uc.tracer.Start(ctx, "EditionsService.ListEditions")
	defer span.End()

	var outEditions []domain.Edition
	outEditions, err = uc.editions.List(ctx, eventID)
	if err != nil {
		return nil, err
	}

	return outEditions, nil
}
