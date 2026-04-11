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

type CommandService struct {
	forms    ports.FormsRepo
	projects ports.ProjectsRepo
	az       *v1.Client
	tx       database.TxRunner
	tracer   trace.Tracer
}

func NewFormCommandService(
	forms ports.FormsRepo,
	projects ports.ProjectsRepo,
	az *v1.Client,
	tx database.TxRunner,
	tracer trace.Tracer,
) *CommandService {
	return &CommandService{
		forms:    forms,
		projects: projects,
		az:       az,
		tx:       tx,
		tracer:   tracer,
	}
}

func (s *CommandService) Create(ctx context.Context, title string, projectID uuid.UUID) (created *contracts.Form, err error) {
	ctx, span := s.tracer.Start(ctx, "FormService.Create")
	defer span.End()

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	var project *contracts.Project
	project, err = s.projects.GetByID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	if err = authz.Require(ctx, s.az,
		authz.Subject("user", sub.ID),
		authz.Permission("create_form"),
		authz.Resource("project", project.ID.String()),
	); err != nil {
		return nil, err
	}

	var form *contracts.Form
	form, err = contracts.NewForm(project.ID, sub.ID, title)
	if err != nil {
		return nil, err
	}

	created, err = s.forms.Create(ctx, *form)
	if err != nil {
		return nil, err
	}

	return created, nil
}
