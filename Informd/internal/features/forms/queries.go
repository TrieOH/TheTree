package forms

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
	forms      ports.FormsRepo
	steps      ports.StepRepo
	namespaces ports.NamespaceRepo
	az         *v1.Client
	tx         database.TxRunner
	tracer     trace.Tracer
}

func NewQueries(
	forms ports.FormsRepo,
	steps ports.StepRepo,
	namespaces ports.NamespaceRepo,
	az *v1.Client,
	tx database.TxRunner,
	tracer trace.Tracer,
) *QueryService {
	return &QueryService{
		forms:      forms,
		steps:      steps,
		namespaces: namespaces,
		az:         az,
		tx:         tx,
		tracer:     tracer,
	}
}

func (s *QueryService) BulkGet(ctx context.Context, ids []uuid.UUID, params contracts.BulkGetParams) (forms []contracts.Form, err error) {
	ctx, span := s.tracer.Start(ctx, "FormService.BulkGet")
	defer span.End()

	_, err = authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	return s.forms.BulkGet(ctx, ids, params)
}
