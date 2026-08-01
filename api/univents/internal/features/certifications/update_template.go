package certifications

import (
	"context"
	"lib/telemetry"
	"univents/internal/authz"
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

	template, err := o.certs.GetTemplateByID(ctx, input.TemplateID)
	if err != nil {
		return nil, err
	}

	edition, err := o.editions.GetByID(ctx, template.EditionID)
	if err != nil {
		return nil, err
	}

	err = authz.Service.CheckEvent(ctx, ident.Sub.ID, edition.EventID, models.EventMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	return o.certs.UpdateTemplate(ctx, input)
}
