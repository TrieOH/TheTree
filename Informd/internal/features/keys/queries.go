package keys

import (
	"Informd/internal/platform/database"
	"Informd/internal/shared/authz"
	"Informd/internal/shared/contracts"
	"Informd/internal/shared/ports"
	"context"

	v1 "github.com/authzed/authzed-go/v1"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

type QueryService struct {
	apiKeys  ports.ApiKeysRepo
	projects ports.ProjectsRepo
	az       *v1.Client
	tx       database.TxRunner
	tracer   trace.Tracer
}

func NewQueries(
	apiKeys ports.ApiKeysRepo,
	projects ports.ProjectsRepo,
	az *v1.Client,
	tx database.TxRunner,
	tracer trace.Tracer,
) *QueryService {
	return &QueryService{
		apiKeys:  apiKeys,
		projects: projects,
		az:       az,
		tx:       tx,
		tracer:   tracer,
	}
}

func (s *QueryService) List(ctx context.Context, projectID uuid.UUID) (ak []contracts.APIKey, err error) {
	ctx, span := s.tracer.Start(ctx, "ApiKeyService.List")
	defer span.End()

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	if err = authz.Require(ctx, s.az,
		authz.Subject("user", sub.ID),
		authz.Permission("list_keys"),
		authz.Resource("project", projectID.String()),
	); err != nil {
		return nil, err
	}

	var keys []contracts.APIKey
	keys, err = s.apiKeys.ListByProject(ctx, projectID)
	if err != nil {
		return nil, err
	}

	return keys, nil
}
