package namespaces

import (
	"Informd/internal/platform/database"
	"Informd/internal/shared/authz"
	"Informd/internal/shared/contracts"
	"Informd/internal/shared/ports"
	"context"

	v1 "github.com/authzed/authzed-go/v1"
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

func (s *QueryService) List(ctx context.Context) (ws []contracts.Namespace, err error) {
	ctx, span := s.tracer.Start(ctx, "NamespaceService.List")
	defer span.End()

	var sub *authz.UserSubject
	if sub, err = authz.RequireSubject(ctx); err != nil {
		return nil, err
	}

	var ids []string
	if ids, err = authz.Lookup(ctx, s.az,
		authz.Subject("user", sub.ID),
		authz.Permission("view"),
		authz.ResourceType("namespace"),
	); err != nil {
		return nil, err
	}

	var namespaces []contracts.Namespace
	if namespaces, err = s.namespaces.ListByIDs(ctx, ids); err != nil {
		return nil, err
	}

	return namespaces, nil
}
