package service

import (
	"context"

	"Informd/models"
	"lib/telemetry"

	"github.com/google/uuid"
)

func (o *Operations) List(ctx context.Context, formID uuid.UUID) ([]models.Step, error) {
	ctx, span := telemetry.StartSpan(ctx, "StepService.List")
	defer span.End()

	_, err := o.Authz.Form(ctx, formID, models.FormMemberRoleMember)
	if err != nil {
		return nil, err
	}

	return o.Steps.List(ctx, formID)
}
