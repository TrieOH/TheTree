package service

import (
	"context"

	"Informd/errs"
	"Informd/models"

	"github.com/google/uuid"
)

func (o *Operations) Archive(ctx context.Context, formID uuid.UUID) (*models.Form, error) {
	form, err := o.Forms.GetByID(ctx, formID)
	if err != nil {
		return nil, err
	}

	if form.Status != models.FormStatusClosed {
		return nil, errs.ErrArchiveNotClosed
	}

	_, err = o.Authz.Form(ctx, form.ID, models.FormMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	return o.Forms.Archive(ctx, formID)
}
