package commands

import (
	"context"
	"lib/telemetry"
	idx "sdk/identityx"

	"Informd/internal/authz"
	"Informd/models"

	"github.com/google/uuid"
)

// TODO: kill this duplicated namespaced route — CheckForm already anchors via the form's namespace.
func (s *Command) DeleteNamespaced(ctx context.Context, _, formID, fieldID uuid.UUID) error {
	ctx, span := telemetry.StartSpan(ctx, "FieldService.DeleteNamespaced")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return err
	}

	err = authz.Service.CheckForm(ctx, ident.Sub.ID, formID, models.FormMemberRoleAdmin)
	if err != nil {
		return err
	}

	return s.fields.Delete(ctx, fieldID)
}
