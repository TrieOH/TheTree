package service

import (
	"context"
	"lib/telemetry"
	idx "sdk/identityx"

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

	err = o.Authz.CheckForm(ctx, ident.Sub.ID, formID, models.FormMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	return o.Fields.UpdateSelectConfig(ctx, payload)
}
