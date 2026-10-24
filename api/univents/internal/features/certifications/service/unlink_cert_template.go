package service

import (
	"context"
	"univents/models"

	"github.com/google/uuid"
)

func (o *Operations) UnlinkCertTemplate(ctx context.Context, templateID, programID uuid.UUID) error {
	template, err := o.Certs.GetTemplateByID(ctx, templateID)
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

	return o.Certs.UnlinkCertTemplate(ctx, templateID, programID)
}
