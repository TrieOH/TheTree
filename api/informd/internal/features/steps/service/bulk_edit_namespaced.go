package service

import (
	"context"

	"Informd/errs"
	"Informd/models"
	"lib/xslices"

	"github.com/google/uuid"
)

// TODO: kill this duplicated namespaced route — authz.Form already anchors via the form's namespace.
func (o *Operations) BulkEditNamespaced(ctx context.Context, formID, _ uuid.UUID, payload []models.UpdateNamespacedFormStepInput) error {
	_, err := o.Authz.Form(ctx, formID, models.FormMemberRoleAdmin)
	if err != nil {
		return err
	}

	for _, p := range payload {
		if p.FormID != formID {
			return errs.ErrStepsAcrossForms
		}
	}

	steps := xslices.MapSlice(payload, models.UpdateNamespacedFormStepInputToStep)
	return o.Steps.BulkEdit(ctx, steps)
}
