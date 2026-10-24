package service

import (
	"context"
	"univents/models"

	"github.com/google/uuid"
)

func (o *Operations) DeleteProgram(ctx context.Context, id uuid.UUID) (*models.Program, error) {
	existing, err := o.Programs.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	edition, err := o.Editions.GetByID(ctx, existing.EditionID)
	if err != nil {
		return nil, err
	}

	_, err = o.Authz.Event(ctx, edition.EventID, models.EventMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	return o.Programs.Delete(ctx, id)
}
