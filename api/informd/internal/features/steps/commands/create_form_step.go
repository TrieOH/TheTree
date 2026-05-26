package commands

import (
	"Informd/models"
	"context"
	"lib/authz"

	"github.com/MintzyG/fun"
)

func (s *Command) Create(ctx context.Context, payload models.CreateFormStepInput) (*models.Step, error) {
	ctx, span := s.tracer.Start(ctx, "StepService.Create")
	defer span.End()

	sub, err := authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	step, err := models.NewStep(payload.FormID, payload.Title, payload.Description, payload.PositionHint)
	if err != nil {
		return nil, err
	}

	member, err := s.forms.GetMember(ctx, sub.ID, payload.FormID)
	if err != nil {
		return nil, err
	}

	if member.Role == models.FormMemberRoleViewer {
		return nil, fun.ErrForbidden("insufficient permissions")
	}

	created, err := s.steps.Create(ctx, *step)
	if err != nil {
		return nil, err
	}

	return created, nil
}
