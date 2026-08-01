package webhook_events

import (
	"context"
	"lib/telemetry"
	"payssage/internal/authz"
	"payssage/models"
	idx "sdk/identityx"

	"github.com/google/uuid"
)

func (o *Operations) GetByID(ctx context.Context, id uuid.UUID) (*models.WebhookEvent, error) {
	ctx, span := telemetry.StartSpan(ctx, "GetEventByID")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	event, err := o.events.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	err = authz.Service.CheckWalletAccess(ctx, ident.Sub.ID, event.WalletID, models.OrganizationRoleMember)
	if err != nil {
		return nil, err
	}

	return event, nil
}
