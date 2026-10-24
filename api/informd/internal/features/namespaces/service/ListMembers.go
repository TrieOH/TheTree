package service

import (
	"context"

	"Informd/models"
	"lib/telemetry"

	"github.com/google/uuid"
)

func (o *Operations) ListMembers(ctx context.Context, namespaceID uuid.UUID) (members []models.NamespaceMember, err error) {
	ctx, span := telemetry.StartSpan(ctx, "NamespaceService.GetMembers")
	defer span.End()

	_, err = o.Authz.Namespace(ctx, namespaceID, models.NamespaceMemberRoleMember)
	if err != nil {
		return nil, err
	}

	members, err = o.Namespaces.ListMembers(ctx, namespaceID)
	if err != nil {
		return nil, err
	}

	return members, nil
}
