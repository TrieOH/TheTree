package commands

import (
	"context"
	idx "sdk/identityx"
	"time"

	"Informd/models"

	"github.com/MintzyG/fun"
)

func (s *Commands) AddMember(ctx context.Context, payload models.AddFormMemberInput) (err error) {
	ctx, span := s.tracer.Start(ctx, "FormService.AddMember")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return err
	}

	if ident.Sub.ID == payload.UserID {
		return fun.ErrBadRequest("users can't add themselves to forms")
	}

	var form *models.Form
	form, err = s.forms.GetByID(ctx, payload.FormID)
	if err != nil {
		return err
	}

	if payload.UserID == form.OwnerID {
		return fun.ErrBadRequest("owner of the form is already a member of the form")
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
	if err == nil {
		return fun.ErrBadRequest("user is already a member of the form")
	}

	newMember := models.FormMember{
		UserID:  payload.UserID,
		FormID:  payload.FormID,
		Role:    payload.Role,
		AddedAt: time.Now(),
		AddedBy: ident.Sub.ID,
	}

	return s.forms.AddMember(ctx, newMember)
}
