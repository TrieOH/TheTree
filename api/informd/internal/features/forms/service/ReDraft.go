package service

import (
	"context"

	"Informd/errs"
	"Informd/models"

	"github.com/google/uuid"
)

func (o *Operations) ReDraft(ctx context.Context, formID uuid.UUID) (*models.Form, error) {
	form, err := o.Forms.GetByID(ctx, formID)
	if err != nil {
		return nil, err
	}

	if form.Status != models.FormStatusOpen {
		return nil, errs.ErrRedraftNotOpen
	}

	count, err := o.Forms.ResponsesCount(ctx, formID)
	if err != nil {
		return nil, err
	}

	if count != 0 {
		return nil, errs.ErrRedraftWithResponses
	}

	_, err = o.Authz.Form(ctx, form.ID, models.FormMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	return o.Forms.ReDraft(ctx, formID)
}
