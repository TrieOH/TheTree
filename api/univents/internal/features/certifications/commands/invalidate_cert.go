package commands

import (
	"context"
	"lib/telemetry"
	"univents/internal/authz"
	"univents/models"

	idx "sdk/identityx"

	"github.com/google/uuid"
)

func (c *Commands) InvalidateCert(ctx context.Context, id uuid.UUID, reason *string) error {
	ctx, span := telemetry.StartSpan(ctx, "CertificationsService.InvalidateCert")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return err
	}

	cert, err := c.certs.GetByID(ctx, id)
	if err != nil {
		return err
	}

	edition, err := c.editions.GetByID(ctx, cert.EditionID)
	if err != nil {
		return err
	}

	err = authz.Service.CheckEvent(ctx, ident.Sub.ID, edition.EventID, models.EventMemberRoleAdmin)
	if err != nil {
		return err
	}

	return c.certs.Invalidate(ctx, id, reason)
}
