package service

import (
	"context"
	"lib/telemetry"
	"univents/models"

	idx "sdk/identityx"
)

func (o *Operations) CreateTemplate(ctx context.Context, input models.CreateCertificationTemplateInput) (*models.CertificationTemplate, error) {
	ctx, span := telemetry.StartSpan(ctx, "CertificationsService.CreateTemplate")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	edition, err := o.Editions.GetByID(ctx, input.EditionID)
	if err != nil {
		return nil, err
	}

	err = o.Authz.CheckEvent(ctx, ident.Sub.ID, edition.EventID, models.EventMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	return o.Certs.CreateTemplate(ctx, input)
}
