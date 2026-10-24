package service

import (
	"context"

	"Informd/models"
	"lib/telemetry"
)

func (o *Operations) Create(ctx context.Context, payload models.CreateFormStepInput) (*models.Step, error) {
	ctx, span := telemetry.StartSpan(ctx, "StepService.Create")
	defer span.End()

	step, err := models.NewStep(payload.FormID, payload.Title, payload.Description, payload.PositionHint)
	if err != nil {
		return nil, err
	}

	_, err = o.Authz.Form(ctx, payload.FormID, models.FormMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	created, err := o.Steps.Create(ctx, *step)
	if err != nil {
		return nil, err
	}

	return created, nil
}
