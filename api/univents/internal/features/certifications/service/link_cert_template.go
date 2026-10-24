package service

import (
	"context"
	"lib/telemetry"
	"univents/models"

	"github.com/google/uuid"
)

func (o *Operations) LinkCertTemplate(ctx context.Context, templateID, programID uuid.UUID) error {
	ctx, span := telemetry.StartSpan(ctx, "CertificationsService.LinkCertTemplate")
	defer span.End()

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

	return o.Certs.LinkCertTemplate(ctx, templateID, programID)
}
