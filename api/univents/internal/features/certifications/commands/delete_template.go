package commands

import (
	"context"
	"lib/telemetry"
	"univents/models"

	idx "sdk/identityx"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

func (c *Commands) DeleteTemplate(ctx context.Context, id uuid.UUID) error {
	ctx, span := telemetry.StartSpan(ctx, "CertificationsService.DeleteTemplate")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return err
	}

	template, err := c.certs.GetTemplateByID(ctx, id)
	if err != nil {
		return err
	}

	edition, err := c.editions.GetByID(ctx, template.EditionID)
	if err != nil {
		return err
	}

	member, err := c.events.GetMember(ctx, edition.EventID, ident.Sub.ID)
	if fun.Is(err, fun.CodeNotFound) {
		return fun.ErrForbidden("insufficient permissions")
	}
	if err != nil {
		return err
	}
	if !member.Role.Minimum(models.EventMemberRoleAdmin) {
		return fun.ErrForbidden("insufficient permissions")
	}

	return c.certs.DeleteTemplate(ctx, id)
}
