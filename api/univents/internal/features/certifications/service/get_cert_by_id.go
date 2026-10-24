package service

import (
	"context"
	"lib/telemetry"
	"univents/models"

	idx "sdk/identityx"

	"github.com/google/uuid"
)

func (o *Operations) GetCertByID(ctx context.Context, id uuid.UUID) (*models.Certification, error) {
	ctx, span := telemetry.StartSpan(ctx, "CertificationsService.GetCertByID")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	cert, err := o.Certs.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if cert.UserID == ident.Sub.ID {
		return cert, nil
	}

	edition, err := o.Editions.GetByID(ctx, cert.EditionID)
	if err != nil {
		return nil, err
	}

	err = o.Authz.CheckEvent(ctx, ident.Sub.ID, edition.EventID, models.EventMemberRoleStaff)
	if err != nil {
		return nil, err
	}

	return cert, nil
}
