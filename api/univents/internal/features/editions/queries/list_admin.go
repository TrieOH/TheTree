package queries

import (
	"context"
	"univents/contracts"

	"github.com/google/uuid"
)

func (uc *Queries) ListAdmin(ctx context.Context, eventID uuid.UUID) (out []contracts.Edition, err error) { // FIXME Pagination
	ctx, span := uc.tracer.Start(ctx, "EditionsService.ListAdmin")
	defer span.End()

	var outEditions []contracts.Edition
	outEditions, err = uc.editions.ListAdmin(ctx, eventID)
	if err != nil {
		return nil, err
	}

	return outEditions, nil
}
