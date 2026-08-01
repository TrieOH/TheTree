package steps

import (
	"context"
	idx "sdk/identityx"

	"Informd/internal/authz"
	"Informd/models"
	"lib/telemetry"
	"lib/xslices"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

func (o *Operations) BulkEdit(ctx context.Context, formID uuid.UUID, payload []models.UpdateFormStepInput) error {
	ctx, span := telemetry.StartSpan(ctx, "StepService.BulkEdit")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return err
	}

	err = authz.Service.CheckForm(ctx, ident.Sub.ID, formID, models.FormMemberRoleAdmin)
	if err != nil {
		return err
	}

	for _, p := range payload {
		if p.FormID != formID {
			return fun.ErrBadRequest("all steps must belong to the same form")
		}
	}

	steps := xslices.MapSlice(payload, models.UpdateFormStepInputToStep)
	return o.steps.BulkEdit(ctx, steps)
}
