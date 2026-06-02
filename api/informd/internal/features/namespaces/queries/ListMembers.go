package queries

import (
	"Informd/models"
	"context"
	"lib/authz"

	"github.com/google/uuid"
)

func (q *Queries) ListMembers(ctx context.Context, namespaceID uuid.UUID) (members []models.NamespaceMember, err error) {
	ctx, span := q.tracer.Start(ctx, "NamespaceService.GetMembers")
	defer span.End()

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	var namespace *models.Namespace
	namespace, err = q.namespaces.GetByID(ctx, namespaceID)
	if err != nil {
		return nil, err
	}

	if sub.ID != namespace.OwnerID {
		_, err = q.namespaces.GetMember(ctx, sub.ID, namespaceID)
		if err != nil {
			return nil, err
		}
	}

	members, err = q.namespaces.ListMembers(ctx, namespaceID)
	if err != nil {
		return nil, err
	}

	return members, nil
}
