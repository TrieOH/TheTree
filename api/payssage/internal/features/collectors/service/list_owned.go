package service

import (
	"context"
	"payssage/models"
	idx "sdk/identityx"
)

func (o *Operations) ListOwned(ctx context.Context) ([]models.Collector, error) {
	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	return o.Collectors.ListByOwner(ctx, ident.Sub.ID)
}
