package namespaces

import (
	"Informd/internal/platform/database"
	"Informd/internal/shared/authz"
	"Informd/internal/shared/contracts"
	"Informd/internal/shared/ports"
	"context"

	"go.opentelemetry.io/otel/trace"
)

type CommandService struct {
	projects ports.NamespaceRepo
	perms    authz.Checker
	tx       database.TxRunner
	tracer   trace.Tracer
}

func NewCommands(
	projects ports.NamespaceRepo,
	perms authz.Checker,
	tx database.TxRunner,
	tracer trace.Tracer,
) *CommandService {
	return &CommandService{
		projects: projects,
		perms:    perms,
		tx:       tx,
		tracer:   tracer,
	}
}

func (s *CommandService) Create(ctx context.Context, name string) (ns *contracts.Namespace, err error) {
	ctx, span := s.tracer.Start(ctx, "NamespaceService.Create")
	defer span.End()

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	var project *contracts.Namespace
	project, err = contracts.NewNamespace(sub.ID, name)
	if err != nil {
		return nil, err
	}

	if err = s.perms.Require(ctx,
		authz.Subject("user", sub.ID),
		authz.Permission("create_namespace"),
		authz.Resource("user", sub.ID.String()),
		map[string]any{"subject_id": sub.ID.String()},
	); err != nil {
		return nil, err
	}

	var created *contracts.Namespace
	created, err = s.projects.Create(ctx, *project)
	if err != nil {
		return nil, err
	}

	if err = s.perms.CreateRelation(ctx,
		"namespace:"+created.ID.String()+"#owner@user:"+sub.ID.String(),
	); err != nil {
		return nil, err
	}

	return created, nil
}
