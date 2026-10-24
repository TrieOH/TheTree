package service

import (
	"context"
	"lib/telemetry"
	idx "sdk/identityx"

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

	err = o.Authz.CheckForm(ctx, ident.Sub.ID, formID, models.FormMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	return o.Fields.UpdateSelectConfig(ctx, payload)
}
