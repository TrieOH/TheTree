package service

import (
	"context"
	idx "sdk/identityx"
	"time"

	"Informd/models"

	"github.com/MintzyG/fun"
)

func (o *Operations) AddMember(ctx context.Context, payload models.AddFormMemberInput) (err error) {
	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return err
	}

	if ident.Sub.ID == payload.UserID {
		return fun.ErrBadRequest("users can't add themselves to forms")
	}

	var form *models.Form
	form, err = o.Forms.GetByID(ctx, payload.FormID)
	if err != nil {
		return err
	}

	if payload.UserID == form.OwnerID {
		return fun.ErrBadRequest("owner of the form is already a member of the form")
	}
	_, err = o.Authz.Form(ctx, form.ID, models.FormMemberRoleAdmin)
	if err != nil {
		return err
	}

	_, err = o.Forms.GetMember(ctx, payload.UserID, form.ID)
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

	return o.Forms.AddMember(ctx, newMember)
}
