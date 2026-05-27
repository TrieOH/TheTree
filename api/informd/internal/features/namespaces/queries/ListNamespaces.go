package queries

import (
	models2 "IdentityX/models"
	"Informd/models"
	"context"
)

func (s *QueryService) ListNamespaces(ctx context.Context) (members []models.Namespace, err error) {
	ctx, span := s.tracer.Start(ctx, "NamespaceService.GetMembers")
	defer span.End()

	var sub *models2.UserSubject
	sub, err = models2.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	ownNamespaces, err := s.namespaces.ListOwned(ctx, sub.ID)
	if err != nil {
		return nil, err
	}

	joinedNamespaces, err := s.namespaces.ListJoined(ctx, sub.ID)
	if err != nil {
		return nil, err
	}

	return append(ownNamespaces, joinedNamespaces...), nil
}
