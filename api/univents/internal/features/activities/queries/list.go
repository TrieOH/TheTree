package queries

import (
	"context"
	"univents/contracts"

	"github.com/google/uuid"
)

func (uc *Queries) List(ctx context.Context, editionID uuid.UUID) (out []contracts.Activity, err error) { // FIXME Pagination
	ctx, span := uc.tracer.Start(ctx, "ActivityService.List")
	defer span.End()

	return uc.activities.List(ctx, editionID)
}
