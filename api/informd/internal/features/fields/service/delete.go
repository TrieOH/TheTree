package service

import (
	"context"
	"lib/telemetry"

	"Informd/models"

	"github.com/google/uuid"
)

func (o *Operations) Delete(ctx context.Context, formID, fieldID uuid.UUID) error {
	ctx, span := telemetry.StartSpan(ctx, "FieldService.Delete")
	defer span.End()

	_, err := o.Authz.Form(ctx, formID, models.FormMemberRoleAdmin)
	if err != nil {
		return err
	}

	return o.Fields.Delete(ctx, fieldID)
}
