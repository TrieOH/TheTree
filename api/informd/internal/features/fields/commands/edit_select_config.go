package commands

import (
	"context"

	"Informd/models"
	"lib/authz"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

func (s *Command) EditSelectConfig(ctx context.Context, formID uuid.UUID, payload models.FieldSelectConfig) (*models.FieldSelectConfig, error) {
	ctx, span := s.tracer.Start(ctx, "FieldService.EditSelectConfig")
	defer span.End()

	sub, err := authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	member, err := s.forms.GetMember(ctx, sub.ID, formID)
	if err != nil {
		return nil, err
	}
	if member.Role == models.FormMemberRoleViewer {
		return nil, fun.ErrForbidden("insufficient permissions")
	}

	return s.fields.UpdateSelectConfig(ctx, payload)
}
