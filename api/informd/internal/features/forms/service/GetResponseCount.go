package service

import (
	"context"
	"uuid"

	"Informd/models"
)

func (o *Operations) GetResponseCount(ctx context.Context, formID uuid.UUID) (int, error) {
	form, err := o.Forms.GetByID(ctx, formID)
	if err != nil {
		return 0, err
	}

	_, err = o.Authz.Form(ctx, form.ID, models.FormMemberRoleMember)
	if err != nil {
		return 0, err
	}

	return o.Forms.ResponsesCount(ctx, formID)
}
