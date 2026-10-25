package service

import (
	"context"

	"github.com/google/uuid"
)

func (o *Operations) Delete(ctx context.Context, id uuid.UUID) error {
	row, err := o.Providers.GetByID(ctx, id)
	if err != nil {
		return err
	}
	err = o.requireProjectAdmin(ctx, row.ProjectID)
	if err != nil {
		return err
	}
	return o.Providers.Delete(ctx, id)
}
