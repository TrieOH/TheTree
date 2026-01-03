package project

import (
	"GoAuth/internal/adapters/observability/tracing"
	"GoAuth/internal/apierr"
	"GoAuth/internal/application/authz"
	"GoAuth/internal/domain/project"
	"GoAuth/internal/ports/inbounds"
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

var _ inbounds.ProjectService = (*UseCase)(nil)

func New(
	projects outbound.ProjectRepository,
) inbounds.ProjectService {
	return &UseCase{projects: projects}
}

// Create handles the business logic for creating a new project.
// It requires a valid principal in the context, generates a new key pair for the project,
// and then creates the project in the database.
func (uc *UseCase) Create(ctx context.Context, in inbounds.CreateProjectInput) (*inbounds.OutputProject, error) {
	ctx, span := usecaseTracer.Start(ctx, "ProjectService.Create")
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

	return inbounds.OutputProjectFromProject(createdProject), nil
}

// GetByID handles the business logic for retrieving a project by its ID.
// It requires a valid principal in the context and that the principal is the owner of the project.
func (uc *UseCase) GetByID(ctx context.Context, projectID string) (*inbounds.OutputProject, error) {
	ctx, span := usecaseTracer.Start(ctx, "ProjectService.GetByID",
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

	return inbounds.OutputProjectFromProject(proj), nil
}

// List handles the business logic for listing all projects for the authenticated user.
func (uc *UseCase) List(ctx context.Context) ([]inbounds.OutputProject, error) {
	ctx, span := usecaseTracer.Start(ctx, "ProjectService.List")
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

	return inbounds.OutputProjectSliceFromProjectSlice(projects), nil
}

// GetJWKS handles the business logic for retrieving the JWKS for a project.
// It retrieves the public key for the project and converts it to a JWK set.
func (uc *UseCase) GetJWKS(ctx context.Context, projectID string) (map[string]any, error) {
	ctx, span := usecaseTracer.Start(ctx, "ProjectService.GetJWKS",
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

// Update handles the business logic for updating a project.
// It requires a valid principal in the context and that the principal is the owner of the project.
// It retrieves the project, updates the fields, and then saves the changes to the database.
func (uc *UseCase) Update(ctx context.Context, in inbounds.UpdateProjectInput) (*inbounds.OutputProject, error) {
	ctx, span := usecaseTracer.Start(ctx, "ProjectService.Update",
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
		ID:          newProject.ID,
		ProjectName: newProject.ProjectName,
		Metadata:    newProject.Metadata,
	},
		principal.UserID,
	)
	if err != nil {
		return nil, err
	}

	return inbounds.OutputProjectFromProject(updatedProject), nil
}

// Delete handles the business logic for deleting a project.
// It requires a valid principal in the context and that the principal is the owner of the project.
func (uc *UseCase) Delete(ctx context.Context, projectID string) error {
	ctx, span := usecaseTracer.Start(ctx, "ProjectService.Delete",
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
