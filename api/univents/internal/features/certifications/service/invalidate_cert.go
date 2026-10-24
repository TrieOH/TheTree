package service

import (
	"context"
	"lib/telemetry"
	"univents/models"

	"github.com/google/uuid"
)

func (o *Operations) InvalidateCert(ctx context.Context, id uuid.UUID, reason *string) error {
	ctx, span := telemetry.StartSpan(ctx, "CertificationsService.InvalidateCert")
	defer span.End()

	cert, err := o.Certs.GetByID(ctx, id)
	if err != nil {
		return err
	}

	edition, err := o.Editions.GetByID(ctx, cert.EditionID)
	if err != nil {
		return err
	}

	_, err = o.Authz.Event(ctx, edition.EventID, models.EventMemberRoleAdmin)
	if err != nil {
		return err
	}

	return o.Certs.Invalidate(ctx, id, reason)
}
