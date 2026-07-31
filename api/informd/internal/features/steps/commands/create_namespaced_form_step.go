package commands

import (
	"context"
	idx "sdk/identityx"

	"Informd/internal/authz"
	"Informd/models"
	"lib/telemetry"
)

// TODO: kill this duplicated namespaced route — CheckForm already anchors via the form's namespace.
func (s *Command) CreateNamespaced(ctx context.Context, payload models.CreateNamespacedFormStepInput) (*models.Step, error) {
	ctx, span := telemetry.StartSpan(ctx, "StepService.CreateNamespaced")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	step, err := models.NewStep(payload.FormID, payload.Title, payload.Description, payload.PositionHint)
	if err != nil {
		return nil, err
	}

	err = authz.Service.CheckForm(ctx, ident.Sub.ID, payload.FormID, models.FormMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	created, err := s.steps.Create(ctx, *step)
	if err != nil {
		return nil, err
	}

	return created, nil
}
