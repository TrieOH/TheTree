package commands

import (
	"Informd/models"
	"context"
	"lib/authz"
)

func (s *CommandService) Create(ctx context.Context, name string) (ns *models.Namespace, err error) {
	ctx, span := s.tracer.Start(ctx, "NamespaceService.Create")
	defer span.End()

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	var project *models.Namespace
	project, err = models.NewNamespace(sub.ID, name)
	if err != nil {
		return nil, err
	}

	var created *models.Namespace
	created, err = s.namespaces.Create(ctx, *project)
	if err != nil {
		return nil, err
	}

	return created, nil
}
