package queries

import (
	"context"
	idx "sdk/identityx"

	"Informd/internal/authz"
	"Informd/models"
	"lib/telemetry"

	"github.com/google/uuid"
)

func (q *Queries) ListMembers(ctx context.Context, namespaceID uuid.UUID) (members []models.NamespaceMember, err error) {
	ctx, span := telemetry.StartSpan(ctx, "NamespaceService.GetMembers")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	err = authz.Service.CheckNamespace(ctx, ident.Sub.ID, namespaceID, models.NamespaceMemberRoleMember)
	if err != nil {
		return nil, err
	}

	members, err = q.namespaces.ListMembers(ctx, namespaceID)
	if err != nil {
		return nil, err
	}

	return members, nil
}
