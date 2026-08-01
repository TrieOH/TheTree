package fields

import (
	"context"
	"lib/telemetry"
	idx "sdk/identityx"

	"Informd/internal/authz"
	"Informd/models"

	"github.com/google/uuid"
)

func (o *Operations) EditSelectConfig(ctx context.Context, formID uuid.UUID, payload models.FieldSelectConfig) (*models.FieldSelectConfig, error) {
	ctx, span := telemetry.StartSpan(ctx, "FieldService.EditSelectConfig")
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
