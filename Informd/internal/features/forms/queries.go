package forms

import (
	"TrieForms/internal/plataform/database"
	"TrieForms/internal/shared/authz"
	"TrieForms/internal/shared/ports"
	"TrieForms/internal/shared/types"
	"context"

	fun "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/TrieOH/goauth-sdk-go"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

type QueryService struct {
	forms    ports.FormsRepo
	projects ports.ProjectsRepo
	gaClient *goauth.Client
	tx       database.TxRunner
	tracer   trace.Tracer
}

func NewFormQueryService(
	forms ports.FormsRepo,
	projects ports.ProjectsRepo,
	gaClient *goauth.Client,
	tx database.TxRunner,
	tracer trace.Tracer,
) *QueryService {
	return &QueryService{
		forms:    forms,
		projects: projects,
		gaClient: gaClient,
		tx:       tx,
		tracer:   tracer,
	}
}

func (s *QueryService) List(ctx context.Context, projectID uuid.UUID) (forms []types.Form, err error) {
	ctx, span := s.tracer.Start(ctx, "FormService.List")
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
		Action("read").
		Scope(project.ScopeID).
		Allowed(ctx)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, fun.NewError("insufficient permissions").Forbidden()
	}

	forms, err = s.forms.ListByProject(ctx, projectID)
	if err != nil {
		return nil, err
	}

	return forms, nil
}
