package queries

import (
	"context"
	idx "sdk/identityx"

	"Informd/internal/authz"
	"Informd/models"
	"lib/telemetry"

	"github.com/google/uuid"
)

func (s *Queries) List(ctx context.Context, formID uuid.UUID) ([]models.Step, error) {
	ctx, span := telemetry.StartSpan(ctx, "StepService.List")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	err = authz.Service.CheckForm(ctx, ident.Sub.ID, formID, models.FormMemberRoleMember)
	if err != nil {
		return nil, err
	}

	return s.steps.List(ctx, formID)
}
