package service

import (
	"context"
	"payssage/models"
	idx "sdk/identityx"
)

func (o *Operations) ListByProfile(ctx context.Context) ([]models.Intent, error) {
	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	return o.Intents.ListByOwner(ctx, ident.Sub.ID)
}
