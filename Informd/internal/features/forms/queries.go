package forms

import (
	"TrieForms/internal/platform/database"
	"TrieForms/internal/shared/authz"
	"TrieForms/internal/shared/contracts"
	"TrieForms/internal/shared/ports"
	"context"

	v1 "github.com/authzed/authzed-go/v1"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

type QueryService struct {
	forms    ports.FormsRepo
	projects ports.ProjectsRepo
	az       *v1.Client
	tx       database.TxRunner
	tracer   trace.Tracer
}

func NewFormQueryService(
	forms ports.FormsRepo,
	projects ports.ProjectsRepo,
	az *v1.Client,
	tx database.TxRunner,
	tracer trace.Tracer,
) *QueryService {
	return &QueryService{
		forms:    forms,
		projects: projects,
		az:       az,
		tx:       tx,
		tracer:   tracer,
	}
}

func (s *QueryService) List(ctx context.Context, projectID uuid.UUID) (forms []contracts.Form, err error) {
	ctx, span := s.tracer.Start(ctx, "FormService.List")
	defer span.End()

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	if err = authz.Require(ctx, s.az,
		authz.Subject("user", sub.ID),
		authz.Permission("list_forms"),
		authz.Resource("project", projectID.String()),
	); err != nil {
		return nil, err
	}

	forms, err = s.forms.ListByProject(ctx, projectID)
	if err != nil {
		return nil, err
	}

	return forms, nil
}
