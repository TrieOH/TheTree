package queries

import (
	"Informd/models"
	"context"
	"lib/authz"
)

func (q *Queries) ListNamespaces(ctx context.Context) (members []models.Namespace, err error) {
	ctx, span := q.tracer.Start(ctx, "NamespaceService.GetMembers")
	defer span.End()

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	ownNamespaces, err := q.namespaces.ListOwned(ctx, sub.ID)
	if err != nil {
		return nil, err
	}

	joinedNamespaces, err := q.namespaces.ListJoined(ctx, sub.ID)
	if err != nil {
		return nil, err
	}

	return append(ownNamespaces, joinedNamespaces...), nil
}
