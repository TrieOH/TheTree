package steps

import (
	"context"
	idx "sdk/identityx"

	"Informd/models"
	"lib/telemetry"
)

func (o *Operations) Create(ctx context.Context, payload models.CreateFormStepInput) (*models.Step, error) {
	ctx, span := telemetry.StartSpan(ctx, "StepService.Create")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	step, err := models.NewStep(payload.FormID, payload.Title, payload.Description, payload.PositionHint)
	if err != nil {
		return nil, err
	}

	err = o.authz.CheckForm(ctx, ident.Sub.ID, payload.FormID, models.FormMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	created, err := o.steps.Create(ctx, *step)
	if err != nil {
		return nil, err
	}

	return created, nil
}
