package forms

import (
	"TrieForms/internal/plataform/database"
	"TrieForms/internal/shared/authz"
	"TrieForms/internal/shared/errx"
	"TrieForms/internal/shared/ports"
	"TrieForms/internal/shared/types"
	"context"
	"encoding/json"

	"github.com/TrieOH/goauth-sdk-go"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

type CommandService struct {
	forms    ports.FormsRepo
	projects ports.ProjectsRepo
	gaClient *goauth.Client
	tx       database.TxRunner
	tracer   trace.Tracer
}

func NewFormCommandService(
	forms ports.FormsRepo,
	projects ports.ProjectsRepo,
	gaClient *goauth.Client,
	tx database.TxRunner,
	tracer trace.Tracer,
) *CommandService {
	return &CommandService{
		forms:    forms,
		projects: projects,
		gaClient: gaClient,
		tx:       tx,
		tracer:   tracer,
	}
}

func (s *CommandService) Create(ctx context.Context, title string, projectID uuid.UUID) (created *types.Form, err error) {
	ctx, span := s.tracer.Start(ctx, "FormService.Create")
	defer span.End()

	ga := s.gaClient

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	var project *types.Project
	project, err = s.projects.GetByID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	var allowed bool
	allowed, err = ga.Authz.Check().User(sub.ID).
		Object("forms").
		Action("create").
		Scope(project.ScopeID).
		Allowed(ctx)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, errx.Forbidden("form").SetMessage("insufficient permissions")
	}

	var form *types.Form
	form, err = types.NewForm(project.ID, sub.ID, title)
	if err != nil {
		return nil, err
	}

	meta := json.RawMessage(`{"color": "#57e389", "icon": "Form", "folder": "forms"}`)
	var scope *goauth.Scope
	var idStr = form.ID.String()
	scope, err = ga.Scopes.CreateWithParent(ctx, project.Name, &idStr, &project.ScopeID, meta)
	if err != nil {
		return nil, err
	}
	form.AddScope(scope.ID)

	created, err = s.forms.Create(ctx, *form)
	if err != nil {
		return nil, err
	}

	return created, nil
}
