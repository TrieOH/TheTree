package intents

import (
	"context"
	"lib/telemetry"
	"payssage/models"
	idx "sdk/identityx"

	"github.com/google/uuid"
)

func (o *Operations) GetByID(ctx context.Context, id uuid.UUID) (*models.Intent, error) {
	ctx, span := telemetry.StartSpan(ctx, "GetByID")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	intent, err := o.intents.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	err = o.authz.CheckWalletAccess(ctx, ident.Sub.ID, intent.WalletID, models.OrganizationRoleMember)
	if err != nil {
		return nil, err
	}

	return intent, nil
}
