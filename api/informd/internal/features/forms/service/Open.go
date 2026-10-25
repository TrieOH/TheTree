package service

import (
	"context"

	"Informd/errs"
	"Informd/models"

	"github.com/google/uuid"
)

func (o *Operations) Open(ctx context.Context, formID uuid.UUID) (*models.Form, error) {
	form, err := o.Forms.GetByID(ctx, formID)
	if err != nil {
		return nil, err
	}

	if form.Status != models.FormStatusDraft {
		return nil, errs.ErrOpenNotDraft
	}

	_, err = o.Authz.Form(ctx, form.ID, models.FormMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	return o.Forms.Open(ctx, formID)
}
