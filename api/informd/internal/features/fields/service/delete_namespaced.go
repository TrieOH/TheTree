package service

import (
	"context"
	"lib/telemetry"

	"Informd/models"

	"github.com/google/uuid"
)

// TODO: kill this duplicated namespaced route — authz.Form already anchors via the form's namespace.
func (o *Operations) DeleteNamespaced(ctx context.Context, _, formID, fieldID uuid.UUID) error {
	ctx, span := telemetry.StartSpan(ctx, "FieldService.DeleteNamespaced")
	defer span.End()

	_, err := o.Authz.Form(ctx, formID, models.FormMemberRoleAdmin)
	if err != nil {
		return err
	}

	return o.Fields.Delete(ctx, fieldID)
}
