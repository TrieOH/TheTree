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
	steps      ports.StepRepo
	namespaces ports.NamespaceRepo
	az         *authzed.Client
	tx         database.TxRunner
	tracer     trace.Tracer
}

func NewCommands(
	forms ports.FormsRepo,
	steps ports.StepRepo,
	namespaces ports.NamespaceRepo,
	az *authzed.Client,
	tx database.TxRunner,
	tracer trace.Tracer,
) *CommandService {
	return &CommandService{
		forms:      forms,
		steps:      steps,
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
			map[string]any{"subject_id": sub.ID.String()},
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

		if err = authz.CreateRelation(ctx, s.az,
			"form:"+created.ID.String()+"#parent_user@user:"+sub.ID.String(),
		); err != nil {
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
		map[string]any{"subject_id": sub.ID.String()},
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

	if err = authz.CreateRelation(ctx, s.az,
		"form:"+created.ID.String()+"#parent_namespace@namespace:"+namespace.ID.String(),
	); err != nil {
		return nil, err
	}

	return created, nil
}

func (s *CommandService) CreateStep(ctx context.Context, formID uuid.UUID, payload CreateStepRequest) (created *contracts.Step, err error) {
	ctx, span := s.tracer.Start(ctx, "FormService.CreateStep")
	defer span.End()

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	_, err = s.forms.GetByID(ctx, formID)
	if err != nil {
		return nil, err
	}

	if err = authz.Require(ctx, s.az,
		authz.Subject("user", sub.ID),
		authz.Permission("edit"),
		authz.Resource("form", formID.String()),
		nil,
	); err != nil {
		return nil, err
	}

	var step *contracts.Step
	step, err = contracts.NewStep(formID, payload.Title, payload.Description, payload.PositionHint)
	if err != nil {
		return nil, err
	}

	created, err = s.steps.Create(ctx, *step)
	if err != nil {
		return nil, err
	}

	return created, nil
}
