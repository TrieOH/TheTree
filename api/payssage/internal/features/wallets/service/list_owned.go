package service

import (
	"context"
	"payssage/models"
	idx "sdk/identityx"
)

func (o *Operations) List(ctx context.Context) ([]models.Wallet, error) {
	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	return o.Wallets.List(ctx, ident.Sub.ID)
}
