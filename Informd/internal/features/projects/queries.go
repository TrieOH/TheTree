package projects

import (
	"TrieForms/internal/plataform/database"
	"TrieForms/internal/shared/authz"
	"TrieForms/internal/shared/ports"
	"TrieForms/internal/shared/types"
	"context"

	fun "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/TrieOH/goauth-sdk-go"
	"go.opentelemetry.io/otel/trace"
)

type QueryService struct {
	projects ports.ProjectsRepo
	gaClient *goauth.Client
	tx       database.TxRunner
	tracer   trace.Tracer
}

func NewProjectQueryService(
	projects ports.ProjectsRepo,
	gaClient *goauth.Client,
	tx database.TxRunner,
	tracer trace.Tracer,
) *QueryService {
	return &QueryService{
		projects: projects,
		gaClient: gaClient,
		tx:       tx,
		tracer:   tracer,
	}
}

func (s *QueryService) List(ctx context.Context) (ws []types.Project, err error) {
	ctx, span := s.tracer.Start(ctx, "ProjectService.List")
	defer span.End()

	ga := s.gaClient

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	var allowed bool
	allowed, err = ga.Authz.Check().User(sub.ID).
		Object("projects").
		Action("read").
		Allowed(ctx)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, fun.NewError("insufficient permissions").Forbidden()
	}

	var projects []types.Project
	projects, err = s.projects.List(ctx, sub.ID)
	if err != nil {
		return nil, err
	}

	return projects, nil
}
