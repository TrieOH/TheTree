package forms

import (
	"Informd/internal/platform/database"
	"Informd/internal/shared/authz"
	"Informd/internal/shared/contracts"
	"Informd/internal/shared/ports"
	"context"

	"github.com/authzed/authzed-go/v1"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

type CommandService struct {
	forms      ports.FormsRepo
	namespaces ports.NamespaceRepo
	az         *authzed.Client
	tx         database.TxRunner
	tracer     trace.Tracer
}

func NewCommands(
	forms ports.FormsRepo,
	namespaces ports.NamespaceRepo,
	az *authzed.Client,
	tx database.TxRunner,
	tracer trace.Tracer,
) *CommandService {
	return &CommandService{
		forms:      forms,
		namespaces: namespaces,
		az:         az,
		tx:         tx,
		tracer:     tracer,
	}
}

func (s *CommandService) Create(ctx context.Context, title string, namespaceID *uuid.UUID) (created *contracts.Form, err error) {
	ctx, span := s.tracer.Start(ctx, "FormService.Create")
	defer span.End()

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	if namespaceID == nil {
		if err = authz.Require(ctx, s.az,
			authz.Subject("user", sub.ID),
			authz.Permission("create_form"),
			authz.Resource("user", sub.ID.String()),
		); err != nil {
			return nil, err
		}

		var form *contracts.Form
		form, err = contracts.NewForm(namespaceID, sub.ID, title)
		if err != nil {
			return nil, err
		}

		created, err = s.forms.Create(ctx, *form)
		if err != nil {
			return nil, err
		}

		return created, nil
	}

	var namespace *contracts.Namespace
	namespace, err = s.namespaces.GetByID(ctx, *namespaceID)
	if err != nil {
		return nil, err
	}

	if err = authz.Require(ctx, s.az,
		authz.Subject("user", sub.ID),
		authz.Permission("create_form"),
		authz.Resource("namespace", namespace.ID.String()),
	); err != nil {
		return nil, err
	}

	var form *contracts.Form
	form, err = contracts.NewForm(&namespace.ID, sub.ID, title)
	if err != nil {
		return nil, err
	}

	created, err = s.forms.Create(ctx, *form)
	if err != nil {
		return nil, err
	}

	return created, nil
}
