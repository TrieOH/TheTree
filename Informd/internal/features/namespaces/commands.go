package namespaces

import (
	"Informd/internal/platform/database"
	"Informd/internal/shared/authz"
	"Informd/internal/shared/contracts"
	"Informd/internal/shared/ports"
	"context"

	v1 "github.com/authzed/authzed-go/v1"
	"go.opentelemetry.io/otel/trace"
)

type CommandService struct {
	projects ports.NamespaceRepo
	az       *v1.Client
	tx       database.TxRunner
	tracer   trace.Tracer
}

func NewCommands(
	projects ports.NamespaceRepo,
	az *v1.Client,
	tx database.TxRunner,
	tracer trace.Tracer,
) *CommandService {
	return &CommandService{
		projects: projects,
		az:       az,
		tx:       tx,
		tracer:   tracer,
	}
}

func (s *CommandService) Create(ctx context.Context, name string) (ws *contracts.Namespace, err error) {
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

	if err = authz.Require(ctx, s.az,
		authz.Subject("user", sub.ID),
		authz.Permission("create_namespace"),
		authz.Resource("platform", "global"),
	); err != nil {
		return nil, err
	}

	var created *contracts.Namespace
	created, err = s.projects.Create(ctx, *project)
	if err != nil {
		return nil, err
	}

	return created, nil
}
