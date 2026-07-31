package commands

import (
	"context"
	"lib/telemetry"
	idx "sdk/identityx"

	"Informd/models"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

func (s *Command) Delete(ctx context.Context, formID, fieldID uuid.UUID) error {
	ctx, span := telemetry.StartSpan(ctx, "FieldService.Delete")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return err
	}

	member, err := s.forms.GetMember(ctx, ident.Sub.ID, formID)
	if err != nil {
		return err
	}
	if member.Role == models.FormMemberRoleViewer {
		return fun.ErrForbidden("insufficient permissions")
	}

	return s.fields.Delete(ctx, fieldID)
}
