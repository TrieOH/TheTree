package service

import (
	"context"
	"lib/telemetry"

	"Informd/models"

	"github.com/google/uuid"
)

func (o *Operations) GetSelectConfig(ctx context.Context, formID, fieldID uuid.UUID) (*models.FieldSelectConfig, error) {
	ctx, span := telemetry.StartSpan(ctx, "FieldService.GetSelectConfig")
	defer span.End()

	_, err := o.Authz.Form(ctx, formID, models.FormMemberRoleMember)
	if err != nil {
		return nil, err
	}

	return o.Fields.GetSelectConfig(ctx, fieldID)
}
