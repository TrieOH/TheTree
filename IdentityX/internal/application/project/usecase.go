package project

import (
	"GoAuth/internal/adapters/observability/tracing"
	"GoAuth/internal/apierr"
	"GoAuth/internal/application/authz"
	"GoAuth/internal/domain/project"
	"GoAuth/internal/ports/outbound"
	"GoAuth/internal/utils"
	"context"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var (
	usecaseTracer = otel.Tracer("auth_usecase")
)

type UseCase struct {
	projects outbound.ProjectRepository
}

func New(
	projects outbound.ProjectRepository,
) *UseCase {
	return &UseCase{projects: projects}
}

func (uc *UseCase) CreateProject(ctx context.Context, in CreateProjectInput) (*OutputProject, error) {
	ctx, span := usecaseTracer.Start(ctx, "ProjectService.CreateProject")
	defer span.End()

	principal, err := authz.RequirePrincipal(ctx)
	if err != nil {
		apierr.RecordDomainError(span, err)
		return nil, err
	}

	tracing.AnnotatePrincipal(span, principal)

	pubKey, privKey, err := utils.GenerateEd25519Keys()
	if err != nil {
		apiErr := apierr.ErrInternal.WithMsg("error generating project keys").WithID(apierr.ProjectErrorGeneratingKeys).WithCause(err)
		apierr.RecordSystemError(span, apiErr)
		return nil, apiErr
	}

	createdProject, err := uc.projects.Create(ctx, project.Project{
		ProjectName: in.ProjectName,
		OwnerID:     principal.UserID,
		Metadata:    in.Metadata,
		IsActive:    true,
		PubKey:      pubKey,
		PrivKey:     []byte(privKey),
	})
	if err != nil {
		return nil, err
	}

	span.SetAttributes(
		attribute.String("project.id", createdProject.ID.String()),
		attribute.String("project.owner_id", createdProject.OwnerID.String()),
		attribute.String("project.name", createdProject.ProjectName),
	)

	return &OutputProject{
		ID:          createdProject.ID,
		ProjectName: createdProject.ProjectName,
		OwnerID:     createdProject.OwnerID,
		Metadata:    createdProject.Metadata,
		IsActive:    createdProject.IsActive,
		CreatedAt:   createdProject.CreatedAt,
		UpdatedAt:   createdProject.UpdatedAt,
	}, nil
}

func (uc *UseCase) GetProjectByID(ctx context.Context, projectID string) (*OutputProject, error) {
	ctx, span := usecaseTracer.Start(ctx, "ProjectService.GetProjectByID",
		trace.WithAttributes(attribute.String("project.id", projectID)),
	)
	defer span.End()

	principal, err := authz.RequirePrincipal(ctx)
	if err != nil {
		apierr.RecordDomainError(span, err)
		return nil, err
	}

	tracing.AnnotatePrincipal(span, principal)

	pid, err := uuid.Parse(projectID)
	if err != nil {
		apiErr := apierr.ErrInvalidInput.WithMsg("invalid project id").WithID(apierr.ProjectInvalidID).WithCause(err)
		apierr.RecordDomainError(span, apiErr)
		return nil, apiErr
	}

	proj, err := uc.projects.GetByID(ctx, pid, principal.UserID)
	if err != nil {
		return nil, err
	}

	span.SetAttributes(
		attribute.String("project.owner_id", proj.OwnerID.String()),
		attribute.String("project.name", proj.ProjectName),
	)

	return &OutputProject{
		ID:          proj.ID,
		ProjectName: proj.ProjectName,
		OwnerID:     proj.OwnerID,
		Metadata:    proj.Metadata,
		IsActive:    proj.IsActive,
		CreatedAt:   proj.CreatedAt,
		UpdatedAt:   proj.UpdatedAt,
	}, nil
}

func (uc *UseCase) ListProjects(ctx context.Context) ([]OutputProject, error) {
	ctx, span := usecaseTracer.Start(ctx, "ProjectService.ListProjects")
	defer span.End()

	principal, err := authz.RequirePrincipal(ctx)
	if err != nil {
		apierr.RecordDomainError(span, err)
		return nil, err
	}

	tracing.AnnotatePrincipal(span, principal)

	projects, err := uc.projects.List(ctx, principal.UserID)
	if err != nil {
		return nil, err
	}

	span.SetAttributes(attribute.Int("projects.count", len(projects)))

	return OutputProjectSliceFromProjectSlice(projects), nil
}

func (uc *UseCase) GetProjectJWKS(ctx context.Context, projectID string) (map[string]any, error) {
	ctx, span := usecaseTracer.Start(ctx, "ProjectService.GetProjectJWKS",
		trace.WithAttributes(attribute.String("project.id", projectID)),
	)
	defer span.End()

	pid, err := uuid.Parse(projectID)
	if err != nil {
		apiErr := apierr.ErrInvalidInput.WithMsg("invalid project id").WithID(apierr.ProjectInvalidID).WithCause(err)
		apierr.RecordDomainError(span, apiErr)
		return nil, apiErr
	}

	pubKey, err := uc.projects.GetPublicKeyByID(ctx, pid)
	if err != nil {
		return nil, err
	}

	parsedKey, err := utils.ParseEd25519PublicKey(pubKey)
	if err != nil {
		apiErr := apierr.ErrInvalidInput.WithMsg("error parsing project public key").WithID(apierr.ProjectErrorParsingKeys).WithCause(err)
		apierr.RecordDomainError(span, apiErr)
		return nil, apiErr
	}

	jwks := utils.PublicKeyToJWK(parsedKey)
	return jwks, nil
}

func (uc *UseCase) UpdateProjectByID(ctx context.Context, in UpdateProjectInput) (*OutputProject, error) {
	ctx, span := usecaseTracer.Start(ctx, "ProjectService.UpdateProjectByID",
		trace.WithAttributes(attribute.String("project.id", in.ProjectID)),
	)
	defer span.End()

	principal, err := authz.RequirePrincipal(ctx)
	if err != nil {
		apierr.RecordDomainError(span, err)
		return nil, err
	}

	tracing.AnnotatePrincipal(span, principal)

	pid, err := uuid.Parse(in.ProjectID)
	if err != nil {
		apiErr := apierr.ErrInvalidInput.WithMsg("invalid project id").WithID(apierr.ProjectInvalidID).WithCause(err)
		apierr.RecordDomainError(span, apiErr)
		return nil, apiErr
	}

	newProject, err := uc.projects.GetByID(ctx, pid, principal.UserID)
	if err != nil {
		return nil, err
	}

	if in.ProjectName != "" {
		newProject.ProjectName = in.ProjectName
	}
	newProject.Metadata = in.Metadata

	updatedProject, err := uc.projects.Update(ctx, project.Project{
		OwnerID:     principal.UserID,
		ID:          newProject.ID,
		ProjectName: newProject.ProjectName,
		Metadata:    newProject.Metadata,
	},
		principal.UserID,
	)
	if err != nil {
		return nil, err
	}

	return OutputProjectFromProject(updatedProject), nil
}

func (uc *UseCase) DeleteProjectByID(ctx context.Context, projectID string) error {
	ctx, span := usecaseTracer.Start(ctx, "ProjectService.DeleteProjectByID",
		trace.WithAttributes(attribute.String("project.id", projectID)),
	)
	defer span.End()

	principal, err := authz.RequirePrincipal(ctx)
	if err != nil {
		apierr.RecordDomainError(span, err)
		return err
	}

	tracing.AnnotatePrincipal(span, principal)

	pid, err := uuid.Parse(projectID)
	if err != nil {
		apiErr := apierr.ErrInvalidInput.WithMsg("invalid project id").WithID(apierr.ProjectInvalidID).WithCause(err)
		apierr.RecordDomainError(span, apiErr)
		return apiErr
	}

	err = uc.projects.Delete(ctx, pid, principal.UserID)
	if err != nil {
		return err
	}

	return nil
}
