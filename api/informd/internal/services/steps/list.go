package steps

import (
	"context"
	idx "sdk/identityx"

	"Informd/models"
	"lib/telemetry"

	"github.com/google/uuid"
)

func (o *Operations) List(ctx context.Context, formID uuid.UUID) ([]models.Step, error) {
	ctx, span := telemetry.StartSpan(ctx, "StepService.List")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	err = o.authz.CheckForm(ctx, ident.Sub.ID, formID, models.FormMemberRoleMember)
	if err != nil {
		return nil, err
	}

	return o.steps.List(ctx, formID)
}
