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
func (o *Operations) EditSelectConfigNamespaced(ctx context.Context, formID, _ uuid.UUID, payload models.FieldSelectConfig) (*models.FieldSelectConfig, error) {
	ctx, span := telemetry.StartSpan(ctx, "FieldService.EditSelectConfigNamespaced")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	err = authz.Service.CheckForm(ctx, ident.Sub.ID, formID, models.FormMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	return o.fields.UpdateSelectConfig(ctx, payload)
}
