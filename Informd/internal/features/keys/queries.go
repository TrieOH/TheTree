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
	apiKeys ports.ApiKeysRepo
	az      *v1.Client
	tx      database.TxRunner
	tracer  trace.Tracer
}

func NewQueries(
	apiKeys ports.ApiKeysRepo,
	az *v1.Client,
	tx database.TxRunner,
	tracer trace.Tracer,
) *QueryService {
	return &QueryService{
		apiKeys: apiKeys,
		az:      az,
		tx:      tx,
		tracer:  tracer,
	}
}

func (s *QueryService) BulkGet(ctx context.Context, ids []uuid.UUID) (keys []contracts.APIKey, err error) {
	ctx, span := s.tracer.Start(ctx, "ApiKeyService.BulkGet")
	defer span.End()

	_, err = authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	return s.apiKeys.BulkGet(ctx, ids)
}
