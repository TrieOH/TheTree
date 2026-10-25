package service

import (
	"context"
	idx "sdk/identityx"

	"Informd/models"

	"github.com/MintzyG/fun"
)

func (o *Operations) RemoveMember(ctx context.Context, payload models.RemoveFormMemberInput) (err error) {
	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return err
	}

	if ident.Sub.ID == payload.UserID {
		return fun.ErrBadRequest("users can't remove themselves from forms")
	}

	var form *models.Form
	form, err = o.Forms.GetByID(ctx, payload.FormID)
	if err != nil {
		return err
	}

	if payload.UserID == form.OwnerID {
		return fun.ErrBadRequest("cannot remove owner of the form")
	}
	_, err = o.Authz.Form(ctx, form.ID, models.FormMemberRoleAdmin)
	if err != nil {
		return err
	}

	_, err = o.Forms.GetMember(ctx, payload.UserID, form.ID)
	if !fun.Is(err, fun.CodeNotFound) {
		return err
	}
	if err != nil {
		return fun.ErrBadRequest("user is not a member of the form")
	}

	return o.Forms.RemoveMember(ctx, payload.UserID, payload.FormID)
}
