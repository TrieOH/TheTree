package commands

import (
	"context"
	idx "sdk/identityx"

	"Informd/models"

	"github.com/MintzyG/fun"
)

func (s *Commands) RemoveMember(ctx context.Context, payload models.RemoveFormMemberInput) (err error) {
	ctx, span := s.tracer.Start(ctx, "FormService.RemoveMember")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return err
	}

	if ident.Sub.ID == payload.UserID {
		return fun.ErrBadRequest("users can't remove themselves from forms")
	}

	var form *models.Form
	if form, err = s.forms.GetByID(ctx, payload.FormID); err != nil {
		return err
	}

	if payload.UserID == form.OwnerID {
		return fun.ErrBadRequest("cannot remove owner of the form")
	}
	if ident.Sub.ID != form.OwnerID {
		member, err := s.forms.GetMember(ctx, ident.Sub.ID, form.ID)
		if err != nil && !fun.Is(err, fun.CodeNotFound) {
			return err
		}
		if err != nil {
			return fun.ErrForbidden("insufficient permissions")
		}
		if member.Role != models.FormMemberRoleAdmin {
			return fun.ErrForbidden("insufficient permissions")
		}
	}

	_, err = s.forms.GetMember(ctx, payload.UserID, form.ID)
	if !fun.Is(err, fun.CodeNotFound) {
		return err
	}
	if err != nil {
		return fun.ErrBadRequest("user is not a member of the form")
	}

	return s.forms.RemoveMember(ctx, payload.UserID, payload.FormID)
}
