package commands

import (
	"context"
	"lib/telemetry"
	"univents/models"

	idx "sdk/identityx"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

func (c *Commands) GetCertByID(ctx context.Context, id uuid.UUID) (*models.Certification, error) {
	ctx, span := telemetry.StartSpan(ctx, "CertificationsService.GetCertByID")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	cert, err := c.certs.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if cert.UserID == ident.Sub.ID {
		return cert, nil
	}

	edition, err := c.editions.GetByID(ctx, cert.EditionID)
	if err != nil {
		return nil, err
	}

	member, err := c.events.GetMember(ctx, edition.EventID, ident.Sub.ID)
	if fun.Is(err, fun.CodeNotFound) {
		return nil, fun.ErrForbidden("insufficient permissions")
	}
	if err != nil {
		return nil, err
	}
	if !member.Role.Minimum(models.EventMemberRoleStaff) {
		return nil, fun.ErrForbidden("insufficient permissions")
	}

	return cert, nil
}
