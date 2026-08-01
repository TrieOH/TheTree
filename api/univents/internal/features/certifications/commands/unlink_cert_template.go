package commands

import (
	"context"
	"lib/telemetry"
	"univents/internal/authz"
	"univents/models"

	idx "sdk/identityx"

	"github.com/google/uuid"
)

func (c *Commands) UnlinkCertTemplate(ctx context.Context, templateID, programID uuid.UUID) error {
	ctx, span := telemetry.StartSpan(ctx, "CertificationsService.UnlinkCertTemplate")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return err
	}

	template, err := c.certs.GetTemplateByID(ctx, templateID)
	if err != nil {
		return err
	}

	edition, err := c.editions.GetByID(ctx, template.EditionID)
	if err != nil {
		return err
	}

	err = authz.Service.CheckEvent(ctx, ident.Sub.ID, edition.EventID, models.EventMemberRoleAdmin)
	if err != nil {
		return err
	}

	return c.certs.UnlinkCertTemplate(ctx, templateID, programID)
}
