package queries

import (
	"context"
	idx "sdk/identityx"

	"Informd/models"
)

func (q *Queries) ListNamespaces(ctx context.Context) (members []models.Namespace, err error) {
	ctx, span := q.tracer.Start(ctx, "NamespaceService.GetMembers")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	ownNamespaces, err := q.namespaces.ListOwned(ctx, ident.Sub.ID)
	if err != nil {
		return nil, err
	}

	joinedNamespaces, err := q.namespaces.ListJoined(ctx, ident.Sub.ID)
	if err != nil {
		return nil, err
	}

	return append(ownNamespaces, joinedNamespaces...), nil
}
