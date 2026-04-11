package projects

import (
	"TrieForms/internal/platform/database"
	"TrieForms/internal/shared/authz"
	"TrieForms/internal/shared/contracts"
	"TrieForms/internal/shared/ports"
	"context"

	v1 "github.com/authzed/authzed-go/v1"
	"go.opentelemetry.io/otel/trace"
)

type CommandService struct {
	projects ports.ProjectsRepo
	az       *v1.Client
	tx       database.TxRunner
	tracer   trace.Tracer
}

func NewProjectCommandService(
	projects ports.ProjectsRepo,
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

func (s *CommandService) Create(ctx context.Context, name string) (ws *contracts.Project, err error) {
	ctx, span := s.tracer.Start(ctx, "ProjectService.Create")
	defer span.End()

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	var project *contracts.Project
	project, err = contracts.NewProject(sub.ID, name)
	if err != nil {
		return nil, err
	}

	if err = authz.Require(ctx, s.az,
		authz.Subject("user", sub.ID),
		authz.Permission("create_project"),
		authz.Resource("platform", "global"),
	); err != nil {
		return nil, err
	}

	var created *contracts.Project
	created, err = s.projects.Create(ctx, *project)
	if err != nil {
		return nil, err
	}

	return created, nil
}
