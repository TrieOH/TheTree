package queries

import (
	"context"
	"lib/telemetry"
	"payssage/internal/authz"
	"payssage/models"
	idx "sdk/identityx"

	"github.com/google/uuid"
)

func (q *Queries) GetByID(ctx context.Context, id uuid.UUID) (*models.WebhookEndpoint, error) {
	ctx, span := telemetry.StartSpan(ctx, "GetEndpointByID")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	endpoint, err := q.endpoints.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	err = authz.Service.CheckWalletAccess(ctx, ident.Sub.ID, endpoint.WalletID, models.OrganizationRoleMember)
	if err != nil {
		return nil, err
	}

	return endpoint, nil
}
