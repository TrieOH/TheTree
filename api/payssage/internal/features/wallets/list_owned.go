package wallets

import (
	"context"
	"lib/telemetry"
	"payssage/models"
	idx "sdk/identityx"
)

func (o *Operations) List(ctx context.Context) ([]models.Wallet, error) {
	ctx, span := telemetry.StartSpan(ctx, "List")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	return o.wallets.List(ctx, ident.Sub.ID)
}
