package fields

import (
	"context"
	"lib/telemetry"
	idx "sdk/identityx"

	"Informd/internal/authz"
	"Informd/models"
	"lib/xslices"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

// TODO: kill this duplicated namespaced route — CheckForm already anchors via the form's namespace.
func (o *Operations) BulkEditNamespaced(ctx context.Context, formID, _ uuid.UUID, payload []models.UpdateNamespacedStepFieldInput) error {
	ctx, span := telemetry.StartSpan(ctx, "FieldService.BulkEditNamespaced")
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
		if p.StepID != payload[0].StepID {
			return fun.ErrBadRequest("all fields must belong to the same step")
		}
	}

	fields := xslices.MapSlice(payload, models.UpdateNamespacedStepFieldInputToField)
	return o.fields.BulkEdit(ctx, fields)
}
