package queries

import (
	"context"
	"univents/contracts"

	"github.com/google/uuid"
)

func (uc *Queries) AdminList(ctx context.Context, editionID uuid.UUID) (out []contracts.Activity, err error) { // FIXME Pagination
	ctx, span := uc.tracer.Start(ctx, "ActivityService.AdminList")
	defer span.End()

	edition, err := uc.editions.GetByID(ctx, editionID)
	if err != nil {
		return nil, err
	}

	return uc.activities.ListAdmin(ctx, edition.ID)
}
