package service

import (
	"context"
	"lib/telemetry"

	"Informd/models"

	"github.com/google/uuid"
)

// TODO: kill this duplicated namespaced route — authz.Form already anchors via the form's namespace.
func (o *Operations) ListNamespaced(ctx context.Context, formID, _, stepID uuid.UUID) ([]models.Field, error) {
	ctx, span := telemetry.StartSpan(ctx, "FieldService.ListNamespaced")
	defer span.End()

	_, err := o.Authz.Form(ctx, formID, models.FormMemberRoleMember)
	if err != nil {
		return nil, err
	}

	return o.Fields.ListByStepID(ctx, stepID)
}
