package service

import (
	"context"
	"lib/telemetry"
	idx "sdk/identityx"

	"Informd/models"

	"github.com/google/uuid"
)

func (o *Operations) List(ctx context.Context, formID, stepID uuid.UUID) ([]models.Field, error) {
	ctx, span := telemetry.StartSpan(ctx, "FieldService.List")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	err = o.Authz.CheckForm(ctx, ident.Sub.ID, formID, models.FormMemberRoleMember)
	if err != nil {
		return nil, err
	}

	return o.Fields.ListByStepID(ctx, stepID)
}
