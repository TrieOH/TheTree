package service

import (
	"GoAuth/internal/apierr"
	"GoAuth/internal/models"
	"GoAuth/internal/utils"
	"context"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func (s *AuthService) CreateProject(ctx context.Context, project models.Project) (*models.Project, error) {
	ctx, span := GoAuthServiceTracer.Start(ctx, "ProjectService.CreateProject")
	defer span.End()

	accessClaims, err := models.GetAccessClaims(ctx)
	if err != nil {
		apierr.RecordDomainError(span, err)
		return nil, err
	}

	annotateAccessClaims(span, accessClaims)

	pubKey, privKey, err := utils.GenerateEd25519Keys()
	if err != nil {
		apiErr := apierr.ErrInternal.WithMsg("error generating project keys").WithID(apierr.ProjectErrorGeneratingKeys).WithCause(err)
		apierr.RecordSystemError(span, apiErr)
		return nil, apiErr
	}

	createdProject, err := s.projectRepo.Create(ctx, models.Project{
		ProjectName: project.ProjectName,
		OwnerID:     accessClaims.Sub.ID,
		Metadata:    project.Metadata,
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

	return createdProject, nil
}

func (s *AuthService) GetProjectByID(ctx context.Context, projectID string) (*models.Project, error) {
	ctx, span := GoAuthServiceTracer.Start(ctx, "ProjectService.GetProjectByID",
		trace.WithAttributes(attribute.String("project.id", projectID)),
	)
	defer span.End()

	accessClaims, err := models.GetAccessClaims(ctx)
	if err != nil {
		apierr.RecordDomainError(span, err)
		return nil, err
	}

	annotateAccessClaims(span, accessClaims)

	pid, err := uuid.Parse(projectID)
	if err != nil {
		apiErr := apierr.ErrInvalidInput.WithMsg("invalid project id").WithID(apierr.ProjectInvalidID).WithCause(err)
		apierr.RecordDomainError(span, apiErr)
		return nil, apiErr
	}

	project, err := s.projectRepo.GetByID(ctx, pid, accessClaims.Sub.ID)
	if err != nil {
		return nil, err
	}

	span.SetAttributes(
		attribute.String("project.owner_id", project.OwnerID.String()),
		attribute.String("project.name", project.ProjectName),
	)

	return project, nil
}

func (s *AuthService) ListProjects(ctx context.Context) ([]models.Project, error) {
	ctx, span := GoAuthServiceTracer.Start(ctx, "ProjectService.ListProjects")
	defer span.End()

	accessClaims, err := models.GetAccessClaims(ctx)
	if err != nil {
		apierr.RecordDomainError(span, err)
		return nil, err
	}

	annotateAccessClaims(span, accessClaims)

	projects, err := s.projectRepo.List(ctx, accessClaims.Sub.ID)
	if err != nil {
		return nil, err
	}

	span.SetAttributes(attribute.Int("projects.count", len(projects)))

	return projects, nil
}

func (s *AuthService) GetProjectKeysByID(ctx context.Context, projectID string) (*models.ProjectKeys, error) {
	ctx, span := GoAuthServiceTracer.Start(ctx, "ProjectService.GetProjectKeysByID",
		trace.WithAttributes(attribute.String("project.id", projectID)),
	)
	defer span.End()

	accessClaims, err := models.GetAccessClaims(ctx)
	if err != nil {
		apierr.RecordDomainError(span, err)
		return nil, err
	}

	annotateAccessClaims(span, accessClaims)

	pid, err := uuid.Parse(projectID)
	if err != nil {
		apiErr := apierr.ErrInvalidInput.WithMsg("invalid project id").WithID(apierr.ProjectInvalidID).WithCause(err)
		apierr.RecordDomainError(span, apiErr)
		return nil, apiErr
	}

	keys, err := s.projectRepo.GetKeysByID(ctx, pid, accessClaims.Sub.ID)
	if err != nil {
		return nil, err
	}

	return keys, nil
}

func (s *AuthService) GetProjectJWKS(ctx context.Context, projectID string) (map[string]any, error) {
	ctx, span := GoAuthServiceTracer.Start(ctx, "ProjectService.GetProjectJWKS",
		trace.WithAttributes(attribute.String("project.id", projectID)),
	)
	defer span.End()

	pid, err := uuid.Parse(projectID)
	if err != nil {
		apiErr := apierr.ErrInvalidInput.WithMsg("invalid project id").WithID(apierr.ProjectInvalidID).WithCause(err)
		apierr.RecordDomainError(span, apiErr)
		return nil, apiErr
	}

	pubKey, err := s.projectRepo.GetPublicKeyByID(ctx, pid)
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

func (s *AuthService) UpdateProjectByID(ctx context.Context, projectID string, project models.Project) (*models.Project, error) {
	ctx, span := GoAuthServiceTracer.Start(ctx, "ProjectService.UpdateProjectByID",
		trace.WithAttributes(attribute.String("project.id", projectID)),
	)
	defer span.End()

	accessClaims, err := models.GetAccessClaims(ctx)
	if err != nil {
		apierr.RecordDomainError(span, err)
		return nil, err
	}

	annotateAccessClaims(span, accessClaims)

	pid, err := uuid.Parse(projectID)
	if err != nil {
		apiErr := apierr.ErrInvalidInput.WithMsg("invalid project id").WithID(apierr.ProjectInvalidID).WithCause(err)
		apierr.RecordDomainError(span, apiErr)
		return nil, apiErr
	}

	newProject, err := s.projectRepo.GetByID(ctx, pid, accessClaims.Sub.ID)
	if err != nil {
		return nil, err
	}

	if project.ProjectName != "" {
		newProject.ProjectName = project.ProjectName
	}
	newProject.Metadata = project.Metadata

	updatedProject, err := s.projectRepo.Update(ctx, models.Project{
		OwnerID:     accessClaims.Sub.ID,
		ID:          newProject.ID,
		ProjectName: newProject.ProjectName,
		Metadata:    newProject.Metadata,
	},
		accessClaims.Sub.ID,
	)
	if err != nil {
		return nil, err
	}

	return updatedProject, nil
}

func (s *AuthService) DeleteProjectByID(ctx context.Context, projectID string) error {
	ctx, span := GoAuthServiceTracer.Start(ctx, "ProjectService.DeleteProjectByID",
		trace.WithAttributes(attribute.String("project.id", projectID)),
	)
	defer span.End()

	accessClaims, err := models.GetAccessClaims(ctx)
	if err != nil {
		apierr.RecordDomainError(span, err)
		return err
	}

	annotateAccessClaims(span, accessClaims)

	pid, err := uuid.Parse(projectID)
	if err != nil {
		apiErr := apierr.ErrInvalidInput.WithMsg("invalid project id").WithID(apierr.ProjectInvalidID).WithCause(err)
		apierr.RecordDomainError(span, apiErr)
		return apiErr
	}

	err = s.projectRepo.Delete(ctx, pid, accessClaims.Sub.ID)
	if err != nil {
		return err
	}

	return nil
}
