package service

import (
	"context"

	"Informd/models"

	"github.com/google/uuid"
)

func (o *Operations) ListMembers(ctx context.Context, formID uuid.UUID) ([]models.FormMember, error) {
	form, err := o.Forms.GetByID(ctx, formID)
	if err != nil {
		return nil, err
	}

	_, err = o.Authz.Form(ctx, form.ID, models.FormMemberRoleMember)
	if err != nil {
		return nil, err
	}

	members, err := o.Forms.ListDirectMembers(ctx, form.ID)
	if err != nil {
		return nil, err
	}

	return members, nil
}
