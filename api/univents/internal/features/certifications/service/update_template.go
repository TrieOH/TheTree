package service

import (
	"context"
	"lib/telemetry"
	"univents/models"
)

func (o *Operations) UpdateTemplate(ctx context.Context, input models.UpdateCertificationTemplateInput) (*models.CertificationTemplate, error) {
	ctx, span := telemetry.StartSpan(ctx, "CertificationsService.UpdateTemplate")
	defer span.End()

	template, err := o.Certs.GetTemplateByID(ctx, input.TemplateID)
	if err != nil {
		return nil, err
	}

	edition, err := o.Editions.GetByID(ctx, template.EditionID)
	if err != nil {
		return nil, err
	}

	_, err = o.Authz.Event(ctx, edition.EventID, models.EventMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	return o.Certs.UpdateTemplate(ctx, input)
}
