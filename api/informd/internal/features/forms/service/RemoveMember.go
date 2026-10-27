package service

import (
	"context"
	"errors"
	idx "sdk/identityx"

	"Informd/errs"
	"Informd/models"
	"lib/database"
)

func (o *Operations) RemoveMember(ctx context.Context, payload models.RemoveFormMemberInput) (err error) {
	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return err
	}

	if ident.Sub.ID == payload.UserID {
		return errs.ErrRemoveSelfFromForm
	}

	var form *models.Form
	form, err = o.Forms.GetByID(ctx, payload.FormID)
	if err != nil {
		return err
	}

	if payload.UserID == form.OwnerID {
		return errs.ErrRemoveFormOwner
	}
	_, err = o.Authz.Form(ctx, form.ID, models.FormMemberRoleAdmin)
	if err != nil {
		return err
	}

	_, err = o.Forms.GetMember(ctx, payload.UserID, form.ID)
	if !errors.Is(err, database.ErrNotFound) {
		return err
	}
	if err != nil {
		return errs.ErrNotFormMember
	}

	return o.Forms.RemoveMember(ctx, payload.UserID, payload.FormID)
}
