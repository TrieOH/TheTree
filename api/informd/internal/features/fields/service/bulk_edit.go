package service

import (
	"context"
	"uuid"

	"Informd/errs"
	"Informd/models"
	"lib/xslices"
)

func (o *Operations) BulkEdit(ctx context.Context, formID uuid.UUID, payload []models.UpdateStepFieldInput) error {
	_, err := o.Authz.Form(ctx, formID, models.FormMemberRoleAdmin)
	if err != nil {
		return err
	}

	for _, p := range payload {
		if p.StepID != payload[0].StepID {
			return errs.ErrFieldsAcrossSteps
		}
	}

	fields := xslices.MapSlice(payload, models.UpdateStepFieldInputToField)
	return o.Fields.BulkEdit(ctx, fields)
}
