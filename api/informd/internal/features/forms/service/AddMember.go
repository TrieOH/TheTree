package service

import (
	"context"
	"errors"
	idx "sdk/identityx"
	"time"

	"Informd/errs"
	"Informd/models"
	"lib/database"
)

func (o *Operations) AddMember(ctx context.Context, payload models.AddFormMemberInput) (err error) {
	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return err
	}

	if ident.Sub.ID == payload.UserID {
		return errs.ErrAddSelfToForm
	}

	var form *models.Form
	form, err = o.Forms.GetByID(ctx, payload.FormID)
	if err != nil {
		return err
	}

	if payload.UserID == form.OwnerID {
		return errs.ErrFormOwnerAlreadyMember
	}
	_, err = o.Authz.Form(ctx, form.ID, models.FormMemberRoleAdmin)
	if err != nil {
		return err
	}

	_, err = o.Forms.GetMember(ctx, payload.UserID, form.ID)
	if !errors.Is(err, database.ErrNotFound) {
		return err
	}
	if err == nil {
		return errs.ErrAlreadyFormMember
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
