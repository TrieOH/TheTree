package commands

import (
	"Informd/models"
	"context"
	"lib/authz"
	"lib/xslices"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

func (s *Command) BulkEdit(ctx context.Context, formID uuid.UUID, payload []models.UpdateFormStepInput) error {
	ctx, span := s.tracer.Start(ctx, "StepService.BulkEdit")
	defer span.End()

	sub, err := authz.RequireSubject(ctx)
	if err != nil {
		return err
	}

	member, err := s.forms.GetMember(ctx, sub.ID, formID)
	if err != nil {
		return err
	}
	if member.Role == models.FormMemberRoleViewer {
		return fun.ErrForbidden("insufficient permissions")
	}

	for _, p := range payload {
		if p.FormID != formID {
			return fun.ErrBadRequest("all steps must belong to the same form")
		}
	}

	steps := xslices.MapSlice(payload, models.UpdateFormStepInputToStep)
	return s.steps.BulkEdit(ctx, steps)
}
