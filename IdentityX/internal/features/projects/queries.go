package projects

import (
	"IdentityX/internal/platform/database"
	"IdentityX/internal/shared/authz"
	"IdentityX/internal/shared/contracts"
	"IdentityX/internal/shared/errx"
	"IdentityX/internal/shared/ports"
	"context"

	"github.com/MintzyG/fail/v3"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type QueryService struct {
	projects ports.ProjectRepository
	users    ports.UserRepository
	logger   *zap.Logger
	tracer   trace.Tracer
	txRunner database.TxRunner
}

func NewQueryService(
	projects ports.ProjectRepository,
	users ports.UserRepository,
	logger *zap.Logger,
	tracer trace.Tracer,
	txRunner database.TxRunner,
) *QueryService {
	return &QueryService{
		projects: projects,
		users:    users,
		logger:   logger,
		tracer:   tracer,
		txRunner: txRunner,
	}
}

// GetByID handles the business logic for retrieving a project by its ID.
// It requires a valid principal in the context and that the principal is the owner of the project.
func (uc *QueryService) GetByID(ctx context.Context, projectID uuid.UUID) (*contracts.Project, error) {
	ctx, span := uc.tracer.Start(ctx, "ProjectService.GetByID",
		trace.WithAttributes(attribute.String("project.id", projectID.String())),
	)
	defer span.End()

	principal, err := authz.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return nil, err
	}

	if principal.ProjectID != nil && *principal.ProjectID != projectID {
		return nil, fail.New(errx.ProjectNotFound).RecordCtx(ctx)
	}

	proj, err := uc.projects.GetByIDExternal(ctx, projectID, principal.UserID)
	if err != nil {
		return nil, err
	}

	span.SetAttributes(
		attribute.String("project.owner_id", proj.OwnerID.String()),
		attribute.String("project.name", proj.ProjectName),
	)

	return proj, nil
}

// List handles the business logic for listing all projects for the authenticated user.
func (uc *QueryService) List(ctx context.Context) ([]contracts.Project, error) {
	ctx, span := uc.tracer.Start(ctx, "ProjectService.List")
	defer span.End()

	principal, err := authz.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return nil, err
	}

	projects, err := uc.projects.List(ctx, principal.UserID)
	if err != nil {
		return nil, err
	}

	span.SetAttributes(attribute.Int("projects.count", len(projects)))

	return projects, nil
}

func (uc *QueryService) ListUsers(ctx context.Context, projectID uuid.UUID) ([]contracts.User, error) {
	ctx, span := uc.tracer.Start(ctx, "ProjectService.ListUsers",
		trace.WithAttributes(attribute.String("project.id", projectID.String())),
	)
	defer span.End()

	users, err := uc.users.ListFromProject(ctx, projectID)
	if err != nil {
		return nil, err
	}

	span.SetAttributes(attribute.Int("users.count", len(users)))

	return users, nil
}

func (uc *QueryService) GetUser(ctx context.Context, projectID, userID uuid.UUID) (*contracts.User, error) {
	ctx, span := uc.tracer.Start(ctx, "ProjectService.GetUser",
		trace.WithAttributes(
			attribute.String("project.id", projectID.String()),
			attribute.String("user.id", userID.String()),
		),
	)
	defer span.End()

	user, err := uc.users.GetByIDFromProject(ctx, userID, projectID)
	if err != nil {
		return nil, err
	}

	return user, nil
}
