package service

import (
	"context"

	"Informd/models"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

func (o *Operations) Archive(ctx context.Context, formID uuid.UUID) (*models.Form, error) {
	form, err := o.Forms.GetByID(ctx, formID)
	if err != nil {
		return nil, err
	}

	if form.Status != models.FormStatusClosed {
		return nil, fun.ErrBadRequest("cannot archive a form not on closed")
	}

	_, err = o.Authz.Form(ctx, form.ID, models.FormMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	return o.Forms.Archive(ctx, formID)
}
