package commands

import (
	"context"

	"Informd/models"
	"lib/authz"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

func (s *Commands) Open(ctx context.Context, formID uuid.UUID) (*models.Form, error) {
	ctx, span := s.tracer.Start(ctx, "FormService.Open")
	defer span.End()

	sub, err := authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	var form *models.Form
	form, err = s.forms.GetByID(ctx, formID)
	if err != nil {
		return nil, err
	}

	if form.Status != models.FormStatusDraft {
		return nil, fun.ErrBadRequest("cannot open a form not on draft")
	}

	if sub.ID != form.OwnerID {
		member, err := s.forms.GetMember(ctx, sub.ID, form.ID)
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

	return s.forms.Open(ctx, formID)
}
