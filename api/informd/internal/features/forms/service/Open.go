package service

import (
	"context"

	"Informd/models"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

func (o *Operations) Open(ctx context.Context, formID uuid.UUID) (*models.Form, error) {
	form, err := o.Forms.GetByID(ctx, formID)
	if err != nil {
		return nil, err
	}

	if form.Status != models.FormStatusDraft {
		return nil, fun.ErrBadRequest("cannot open a form not on draft")
	}

	_, err = o.Authz.Form(ctx, form.ID, models.FormMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	return o.Forms.Open(ctx, formID)
}
