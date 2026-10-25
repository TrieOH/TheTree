package service

import (
	"context"

	"Informd/models"

	"github.com/google/uuid"
)

func (o *Operations) ListMembers(ctx context.Context, namespaceID uuid.UUID) (members []models.NamespaceMember, err error) {
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
