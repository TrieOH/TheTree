package service

import (
	"context"
	"lib/telemetry"
	idx "sdk/identityx"

	"Informd/models"

	"github.com/google/uuid"
)

func (o *Operations) GetSelectConfig(ctx context.Context, formID, fieldID uuid.UUID) (*models.FieldSelectConfig, error) {
	ctx, span := telemetry.StartSpan(ctx, "FieldService.GetSelectConfig")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	err = o.Authz.CheckForm(ctx, ident.Sub.ID, formID, models.FormMemberRoleMember)
	if err != nil {
		return nil, err
	}

	return o.Fields.GetSelectConfig(ctx, fieldID)
}
