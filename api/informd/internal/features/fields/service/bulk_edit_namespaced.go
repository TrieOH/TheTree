package service

import (
	"context"

	"Informd/models"
	"lib/xslices"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

// TODO: kill this duplicated namespaced route — authz.Form already anchors via the form's namespace.
func (o *Operations) BulkEditNamespaced(ctx context.Context, formID, _ uuid.UUID, payload []models.UpdateNamespacedStepFieldInput) error {
	_, err := o.Authz.Form(ctx, formID, models.FormMemberRoleAdmin)
	if err != nil {
		return err
	}

	for _, p := range payload {
		if p.StepID != payload[0].StepID {
			return fun.ErrBadRequest("all fields must belong to the same step")
		}
	}

	fields := xslices.MapSlice(payload, models.UpdateNamespacedStepFieldInputToField)
	return o.Fields.BulkEdit(ctx, fields)
}
