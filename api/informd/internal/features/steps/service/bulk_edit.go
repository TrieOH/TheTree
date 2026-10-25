package service

import (
	"context"

	"Informd/errs"
	"Informd/models"
	"lib/xslices"

	"github.com/google/uuid"
)

func (o *Operations) BulkEdit(ctx context.Context, formID uuid.UUID, payload []models.UpdateFormStepInput) error {
	_, err := o.Authz.Form(ctx, formID, models.FormMemberRoleAdmin)
	if err != nil {
		return err
	}

	for _, p := range payload {
		if p.FormID != formID {
			return errs.ErrStepsAcrossForms
		}
	}

	steps := xslices.MapSlice(payload, models.UpdateFormStepInputToStep)
	return o.Steps.BulkEdit(ctx, steps)
}
