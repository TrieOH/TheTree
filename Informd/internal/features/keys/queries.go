package keys

import (
	"TrieForms/internal/plataform/database"
	"TrieForms/internal/shared/authz"
	"TrieForms/internal/shared/errx"
	"TrieForms/internal/shared/ports"
	"TrieForms/internal/shared/types"
	"context"

	"github.com/TrieOH/goauth-sdk-go"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

type QueryService struct {
	apiKeys  ports.ApiKeysRepo
	projects ports.ProjectsRepo
	gaClient *goauth.Client
	tx       database.TxRunner
	tracer   trace.Tracer
}

func NewApiKeyQueryService(
	apiKeys ports.ApiKeysRepo,
	projects ports.ProjectsRepo,
	gaClient *goauth.Client,
	tx database.TxRunner,
	tracer trace.Tracer,
) *QueryService {
	return &QueryService{
		apiKeys:  apiKeys,
		projects: projects,
		gaClient: gaClient,
		tx:       tx,
		tracer:   tracer,
	}
}

func (s *QueryService) List(ctx context.Context, projectID uuid.UUID) (ak []types.APIKey, err error) {
	ctx, span := s.tracer.Start(ctx, "CommandService.List")
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
		Object("api_keys").
		Action("read").
		Scope(project.ScopeID).
		Allowed(ctx)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, errx.Forbidden("api key").SetMessage("insufficient permissions")
	}

	var keys []types.APIKey
	keys, err = s.apiKeys.ListByProject(ctx, projectID)
	if err != nil {
		return nil, err
	}

	return keys, nil
}
