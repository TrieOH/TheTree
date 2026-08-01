package namespaces

import (
	"context"
	idx "sdk/identityx"

	"Informd/models"
	"lib/telemetry"
)

func (o *Operations) ListNamespaces(ctx context.Context) (members []models.Namespace, err error) {
	ctx, span := telemetry.StartSpan(ctx, "NamespaceService.GetMembers")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	ownNamespaces, err := o.namespaces.ListOwned(ctx, ident.Sub.ID)
	if err != nil {
		return nil, err
	}

	joinedNamespaces, err := o.namespaces.ListJoined(ctx, ident.Sub.ID)
	if err != nil {
		return nil, err
	}

	return append(ownNamespaces, joinedNamespaces...), nil
}
