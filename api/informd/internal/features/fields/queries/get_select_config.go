package queries

import (
	"context"
	"lib/telemetry"
	idx "sdk/identityx"

	"Informd/internal/authz"
	"Informd/models"

	"github.com/google/uuid"
)

func (s *Queries) GetSelectConfig(ctx context.Context, formID, fieldID uuid.UUID) (*models.FieldSelectConfig, error) {
	ctx, span := telemetry.StartSpan(ctx, "FieldService.GetSelectConfig")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	err = authz.Service.CheckForm(ctx, ident.Sub.ID, formID, models.FormMemberRoleMember)
	if err != nil {
		return nil, err
	}

	return s.fields.GetSelectConfig(ctx, fieldID)
}
