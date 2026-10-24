package service

import (
	"context"
	"lib/telemetry"
	"univents/models"

	idx "sdk/identityx"
)

func (o *Operations) UpdateTemplate(ctx context.Context, input models.UpdateCertificationTemplateInput) (*models.CertificationTemplate, error) {
	ctx, span := telemetry.StartSpan(ctx, "CertificationsService.UpdateTemplate")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	template, err := o.Certs.GetTemplateByID(ctx, input.TemplateID)
	if err != nil {
		return nil, err
	}

	edition, err := o.Editions.GetByID(ctx, template.EditionID)
	if err != nil {
		return nil, err
	}

	err = o.Authz.CheckEvent(ctx, ident.Sub.ID, edition.EventID, models.EventMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	return o.Certs.UpdateTemplate(ctx, input)
}
