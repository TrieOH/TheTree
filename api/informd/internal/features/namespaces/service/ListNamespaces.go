package service

import (
	"context"
	idx "sdk/identityx"

	"Informd/models"
)

func (o *Operations) ListNamespaces(ctx context.Context) (members []models.Namespace, err error) {
	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	ownNamespaces, err := o.Namespaces.ListOwned(ctx, ident.Sub.ID)
	if err != nil {
		return nil, err
	}

	joinedNamespaces, err := o.Namespaces.ListJoined(ctx, ident.Sub.ID)
	if err != nil {
		return nil, err
	}

	return append(ownNamespaces, joinedNamespaces...), nil
}
