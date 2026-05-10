package namespaces

import (
	"Informd/contracts"
	"Informd/internal/shared/ports"
	"context"
	"lib/authz"
	"lib/database"

	v1 "github.com/authzed/authzed-go/v1"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

type QueryService struct {
	namespaces ports.NamespaceRepo
	az         *v1.Client
	tx         database.TxRunner
	tracer     trace.Tracer
}

func NewQueries(
	namespaces ports.NamespaceRepo,
	az *v1.Client,
	tx database.TxRunner,
	tracer trace.Tracer,
) *QueryService {
	return &QueryService{
		namespaces: namespaces,
		az:         az,
		tx:         tx,
		tracer:     tracer,
	}
}

func (s *QueryService) BulkGet(ctx context.Context, ids []uuid.UUID) (ns []contracts.Namespace, err error) {
	ctx, span := s.tracer.Start(ctx, "NamespaceService.BulkGet")
	defer span.End()

	_, err = authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	return s.namespaces.BulkGet(ctx, ids)
}
