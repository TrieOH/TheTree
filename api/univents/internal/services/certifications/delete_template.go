package certifications

import (
	"context"
	"lib/telemetry"
	"univents/models"

	idx "sdk/identityx"

	"github.com/google/uuid"
)

func (o *Operations) DeleteTemplate(ctx context.Context, id uuid.UUID) error {
	ctx, span := telemetry.StartSpan(ctx, "CertificationsService.DeleteTemplate")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return err
	}

	template, err := o.certs.GetTemplateByID(ctx, id)
	if err != nil {
		return err
	}

	edition, err := o.editions.GetByID(ctx, template.EditionID)
	if err != nil {
		return err
	}

	err = o.authz.CheckEvent(ctx, ident.Sub.ID, edition.EventID, models.EventMemberRoleAdmin)
	if err != nil {
		return err
	}

	return o.certs.DeleteTemplate(ctx, id)
}
