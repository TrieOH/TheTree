package commands

import (
	"context"

	"Informd/models"
	"lib/authz"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

func (s *Command) Delete(ctx context.Context, formID, fieldID uuid.UUID) error {
	ctx, span := s.tracer.Start(ctx, "FieldService.Delete")
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

	return s.fields.Delete(ctx, fieldID)
}
