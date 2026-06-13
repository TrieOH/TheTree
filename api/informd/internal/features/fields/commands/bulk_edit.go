package commands

import (
	"context"

	"Informd/models"
	"lib/authz"
	"lib/xslices"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

func (s *Command) BulkEdit(ctx context.Context, formID uuid.UUID, payload []models.UpdateStepFieldInput) error {
	ctx, span := s.tracer.Start(ctx, "FieldService.BulkEdit")
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
		if p.StepID != payload[0].StepID {
			return fun.ErrBadRequest("all fields must belong to the same step")
		}
	}

	fields := xslices.MapSlice(payload, models.UpdateStepFieldInputToField)
	return s.fields.BulkEdit(ctx, fields)
}
