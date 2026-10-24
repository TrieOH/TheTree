package service

import (
	"context"
	"univents/models"

	"github.com/google/uuid"
)

func (o *Operations) DeleteTemplate(ctx context.Context, templateID uuid.UUID) error {
	template, err := o.Repo.GetByID(ctx, templateID)
	if err != nil {
		return err
	}

	edition, err := o.Editions.GetByID(ctx, template.EditionID)
	if err != nil {
		return err
	}

	_, err = o.Authz.Event(ctx, edition.EventID, models.EventMemberRoleAdmin)
	if err != nil {
		return err
	}

	return o.Repo.Delete(ctx, templateID)
}
