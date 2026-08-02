package oauth_providers

import (
	"context"
	"lib/telemetry"

	"github.com/google/uuid"
)

func (o *Operations) Delete(ctx context.Context, id uuid.UUID) error {
	ctx, span := telemetry.StartSpan(ctx, "Delete")
	defer span.End()

	row, err := o.providers.GetByID(ctx, id)
	if err != nil {
		return err
	}
	err = o.requireProjectAdmin(ctx, row.ProjectID)
	if err != nil {
		return err
	}
	return o.providers.Delete(ctx, id)
}
