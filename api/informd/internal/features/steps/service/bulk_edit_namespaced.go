package service

import (
	"context"

	"Informd/models"
	"lib/telemetry"
	"lib/xslices"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

// TODO: kill this duplicated namespaced route — authz.Form already anchors via the form's namespace.
func (o *Operations) BulkEditNamespaced(ctx context.Context, formID, _ uuid.UUID, payload []models.UpdateNamespacedFormStepInput) error {
	ctx, span := telemetry.StartSpan(ctx, "StepService.BulkEditNamespaced")
	defer span.End()

	_, err := o.Authz.Form(ctx, formID, models.FormMemberRoleAdmin)
	if err != nil {
		return err
	}

	for _, p := range payload {
		if p.FormID != formID {
			return fun.ErrBadRequest("all steps must belong to the same form")
		}
	}

	steps := xslices.MapSlice(payload, models.UpdateNamespacedFormStepInputToStep)
	return o.Steps.BulkEdit(ctx, steps)
}
