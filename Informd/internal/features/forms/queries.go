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
	namespaces ports.NamespaceRepo
	az         *v1.Client
	tx         database.TxRunner
	tracer     trace.Tracer
}

func NewQueries(
	forms ports.FormsRepo,
	namespaces ports.NamespaceRepo,
	az *v1.Client,
	tx database.TxRunner,
	tracer trace.Tracer,
) *QueryService {
	return &QueryService{
		forms:      forms,
		namespaces: namespaces,
		az:         az,
		tx:         tx,
		tracer:     tracer,
	}
}

func (s *QueryService) List(ctx context.Context, namespaceID *uuid.UUID) (forms []contracts.Form, err error) {
	ctx, span := s.tracer.Start(ctx, "FormService.List")
	defer span.End()

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	if namespaceID == nil {
		if err = authz.Require(ctx, s.az,
			authz.Subject("user", sub.ID),
			authz.Permission("list_forms"),
			authz.Resource("user", sub.ID.String()),
		); err != nil {
			return nil, err
		}

		forms, err = s.forms.List(ctx, sub.ID)
		if err != nil {
			return nil, err
		}

		return forms, nil
	}

	if err = authz.Require(ctx, s.az,
		authz.Subject("user", sub.ID),
		authz.Permission("list_forms"),
		authz.Resource("namespace", namespaceID.String()),
	); err != nil {
		return nil, err
	}

	forms, err = s.forms.ListByNamespace(ctx, namespaceID)
	if err != nil {
		return nil, err
	}

	return forms, nil
}
