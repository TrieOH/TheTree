package service

import (
	"context"
	"lib/telemetry"
	"univents/models"

	"github.com/google/uuid"
)

func (o *Operations) DeleteTemplate(ctx context.Context, id uuid.UUID) error {
	ctx, span := telemetry.StartSpan(ctx, "CertificationsService.DeleteTemplate")
	defer span.End()

	template, err := o.Certs.GetTemplateByID(ctx, id)
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

	return o.Certs.DeleteTemplate(ctx, id)
}
