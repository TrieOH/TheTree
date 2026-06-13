package commands

import (
	"context"
	"time"

	"Informd/models"
	"lib/authz"

	"github.com/MintzyG/fun"
)

func (s *Commands) AddMember(ctx context.Context, payload models.AddFormMemberInput) (err error) {
	ctx, span := s.tracer.Start(ctx, "FormService.AddMember")
	defer span.End()

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return err
	}

	if sub.ID == payload.UserID {
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
	if sub.ID != form.OwnerID {
		member, err := s.forms.GetMember(ctx, sub.ID, form.ID)
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
		AddedBy: sub.ID,
	}

	return s.forms.AddMember(ctx, newMember)
}
