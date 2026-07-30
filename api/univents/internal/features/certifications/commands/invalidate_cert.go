package commands

import (
	"context"
	"lib/telemetry"
	"univents/models"

	idx "sdk/identityx"

	"github.com/MintzyG/fun"
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

	return c.certs.Invalidate(ctx, id, reason)
}
