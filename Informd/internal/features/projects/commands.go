package projects

import (
	"TrieForms/internal/plataform/database"
	"TrieForms/internal/shared/authz"
	"TrieForms/internal/shared/errx"
	"TrieForms/internal/shared/ports"
	"TrieForms/internal/shared/types"
	"context"
	"encoding/json"

	"github.com/TrieOH/goauth-sdk-go"
	"go.opentelemetry.io/otel/trace"
)

type CommandService struct {
	projects ports.ProjectsRepo
	gaClient *goauth.Client
	tx       database.TxRunner
	tracer   trace.Tracer
}

func NewProjectCommandService(
	projects ports.ProjectsRepo,
	gaClient *goauth.Client,
	tx database.TxRunner,
	tracer trace.Tracer,
) *CommandService {
	return &CommandService{
		projects: projects,
		gaClient: gaClient,
		tx:       tx,
		tracer:   tracer,
	}
}

func (s *CommandService) Create(ctx context.Context, name string) (ws *types.Project, err error) {
	ctx, span := s.tracer.Start(ctx, "ProjectService.Create")
	defer span.End()

	ga := s.gaClient

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	var project *types.Project
	project, err = types.NewProject(sub.ID, name)
	if err != nil {
		return nil, err
	}

	var allowed bool
	allowed, err = ga.Authz.Check().User(sub.ID).
		Object("projects").
		Action("create").
		Allowed(ctx)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, errx.Forbidden("project").SetMessage("insufficient permissions")
	}

	meta := json.RawMessage(`{"color": "#6a07e3", "icon": "Shield"}`)
	var scope *goauth.Scope
	var idStr = project.ID.String()
	scope, err = ga.Scopes.CreateWithParent(ctx, project.Name, &idStr, nil, meta)
	if err != nil {
		return nil, err
	}
	project.AddScope(scope.ID)

	var created *types.Project
	created, err = s.projects.Create(ctx, *project)
	if err != nil {
		return nil, err
	}

	return created, nil
}
