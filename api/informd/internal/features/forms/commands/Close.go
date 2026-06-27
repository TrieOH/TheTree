package commands

import (
	"context"
	idx "sdk/identityx"

	"Informd/models"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

func (s *Commands) Close(ctx context.Context, formID uuid.UUID) (*models.Form, error) {
	ctx, span := s.tracer.Start(ctx, "FormService.Close")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	var form *models.Form
	form, err = s.forms.GetByID(ctx, formID)
	if err != nil {
		return nil, err
	}

	if form.Status != models.FormStatusOpen {
		return nil, fun.ErrBadRequest("cannot close a form not on open")
	}

	if ident.Sub.ID != form.OwnerID {
		member, err := s.forms.GetMember(ctx, ident.Sub.ID, form.ID)
		if err != nil && !fun.Is(err, fun.CodeNotFound) {
			return nil, err
		}
		if err != nil {
			return nil, fun.ErrForbidden("insufficient permissions")
		}
		if member.Role != models.FormMemberRoleAdmin {
			return nil, fun.ErrForbidden("insufficient permissions")
		}
	}

	return s.forms.Close(ctx, formID)
}
