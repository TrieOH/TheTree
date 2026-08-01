package certifications

import (
	"context"
	"lib/telemetry"
	"univents/internal/authz"
	"univents/models"

	idx "sdk/identityx"

	"github.com/google/uuid"
)

func (o *Operations) LinkCertTemplate(ctx context.Context, templateID, programID uuid.UUID) error {
	ctx, span := telemetry.StartSpan(ctx, "CertificationsService.LinkCertTemplate")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return err
	}

	template, err := o.certs.GetTemplateByID(ctx, templateID)
	if err != nil {
		return err
	}

	edition, err := o.editions.GetByID(ctx, template.EditionID)
	if err != nil {
		return err
	}

	err = authz.Service.CheckEvent(ctx, ident.Sub.ID, edition.EventID, models.EventMemberRoleAdmin)
	if err != nil {
		return err
	}

	return o.certs.LinkCertTemplate(ctx, templateID, programID)
}
