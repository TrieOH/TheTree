package steps

import (
	"context"
	idx "sdk/identityx"

	"Informd/models"
	"lib/telemetry"
	"lib/xslices"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

// TODO: kill this duplicated namespaced route — CheckForm already anchors via the form's namespace.
func (o *Operations) BulkEditNamespaced(ctx context.Context, formID, _ uuid.UUID, payload []models.UpdateNamespacedFormStepInput) error {
	ctx, span := telemetry.StartSpan(ctx, "StepService.BulkEditNamespaced")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return err
	}

	err = o.authz.CheckForm(ctx, ident.Sub.ID, formID, models.FormMemberRoleAdmin)
	if err != nil {
		return err
	}

	for _, p := range payload {
		if p.FormID != formID {
			return fun.ErrBadRequest("all steps must belong to the same form")
		}
	}

	steps := xslices.MapSlice(payload, models.UpdateNamespacedFormStepInputToStep)
	return o.steps.BulkEdit(ctx, steps)
}
