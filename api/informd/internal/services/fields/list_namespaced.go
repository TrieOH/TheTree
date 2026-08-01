package fields

import (
	"context"
	"lib/telemetry"
	idx "sdk/identityx"

	"Informd/internal/authz"
	"Informd/models"

	"github.com/google/uuid"
)

// TODO: kill this duplicated namespaced route — CheckForm already anchors via the form's namespace.
func (o *Operations) ListNamespaced(ctx context.Context, formID, _, stepID uuid.UUID) ([]models.Field, error) {
	ctx, span := telemetry.StartSpan(ctx, "FieldService.ListNamespaced")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	err = authz.Service.CheckForm(ctx, ident.Sub.ID, formID, models.FormMemberRoleMember)
	if err != nil {
		return nil, err
	}

	return o.fields.ListByStepID(ctx, stepID)
}
