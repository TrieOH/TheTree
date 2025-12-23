package service

import (
	"GoAuth/internal/models"
	"GoAuth/internal/utils"
	"context"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/google/uuid"
)

func (s *AuthService) CreateProject(ctx context.Context, project models.Project) (*models.Project, *resp.Response) {
	accessClaims, err := models.GetAccessClaims(ctx)
	if err != nil {
		return nil, resp.InternalServerError().AddTrace(err)
	}

	pubKey, privKey, err := utils.GenerateEd25519Keys()
	if err != nil {
		return nil, resp.InternalServerError("error generating public and private project keys").AddTrace(err)
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
		return nil, resp.InternalServerError("error creating project").WithTracePrefix("database-error").AddTrace(err)
	}

	return createdProject, nil
}

func (s *AuthService) GetProjectByID(ctx context.Context, projectID string) (*models.Project, *resp.Response) {
	accessClaims, err := models.GetAccessClaims(ctx)
	if err != nil {
		return nil, resp.InternalServerError().AddTrace(err)
	}

	pid, err := uuid.Parse(projectID)
	if err != nil {
		return nil, resp.InternalServerError("error parsing project uuid").AddTrace(err)
	}

	project, err := s.projectRepo.GetByID(ctx, pid, accessClaims.Sub.ID)

	if err != nil {
		return nil, resp.InternalServerError("error fetching project").WithTracePrefix("database-error").AddTrace(err)
	}

	return project, nil
}

func (s *AuthService) ListProjects(ctx context.Context) ([]models.Project, *resp.Response) {
	accessClaims, err := models.GetAccessClaims(ctx)
	if err != nil {
		return nil, resp.InternalServerError().AddTrace(err)
	}

	projects, err := s.projectRepo.List(ctx, accessClaims.Sub.ID)

	if err != nil {
		return nil, resp.InternalServerError("error fetching projects").WithTracePrefix("database-error").AddTrace(err)
	}

	return projects, nil
}

func (s *AuthService) GetProjectKeysByID(ctx context.Context, projectId string) (*models.ProjectKeys, *resp.Response) {
	accessClaims, err := models.GetAccessClaims(ctx)
	if err != nil {
		return nil, resp.InternalServerError().AddTrace(err)
	}

	pid, err := uuid.Parse(projectId)
	if err != nil {
		return nil, resp.BadRequest("error parsing project id").AddTrace(err)
	}

	keys, err := s.projectRepo.GetKeysByID(ctx, pid, accessClaims.Sub.ID)

	if err != nil {
		return nil, resp.InternalServerError("error fetching project keys").WithTracePrefix("database-error").AddTrace(err)
	}

	return keys, nil
}

func (s *AuthService) GetProjectJWKS(ctx context.Context, projectId string) (map[string]any, *resp.Response) {
	pid, err := uuid.Parse(projectId)
	if err != nil {
		return nil, resp.BadRequest("error parsing project id").AddTrace(err)
	}

	pubKey, err := s.projectRepo.GetPublicKeyByID(ctx, pid)
	if err != nil {
		return nil, resp.InternalServerError("error fetching project public key").WithTracePrefix("database-error").AddTrace(err)
	}

	parsedKey, err := utils.ParseEd25519PublicKey(pubKey)
	if err != nil {
		return nil, resp.BadRequest("error parsing project public key").AddTrace(err)
	}

	jwks := utils.PublicKeyToJWK(parsedKey)

	return jwks, nil
}

func (s *AuthService) UpdateProjectByID(ctx context.Context, ProjectID string, project models.Project) (*models.Project, *resp.Response) {
	accessClaims, err := models.GetAccessClaims(ctx)
	if err != nil {
		return nil, resp.InternalServerError().AddTrace(err)
	}

	pid, err := uuid.Parse(ProjectID)
	if err != nil {
		return nil, resp.BadRequest("error parsing project id").AddTrace(err)
	}

	newProject, err := s.projectRepo.GetByID(ctx, pid, accessClaims.Sub.ID)
	if err != nil {
		return nil, resp.InternalServerError("error fetching project").AddTrace(err)
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
		return nil, resp.InternalServerError("error updating project").WithTracePrefix("database-error").AddTrace(err)
	}

	return updatedProject, nil
}

func (s *AuthService) DeleteProjectByID(ctx context.Context, projectId string) *resp.Response {
	accessClaims, err := models.GetAccessClaims(ctx)
	if err != nil {
		return resp.InternalServerError().AddTrace(err)
	}

	pid, err := uuid.Parse(projectId)
	if err != nil {
		return resp.BadRequest("error parsing project id").AddTrace(err)
	}

	err = s.projectRepo.Delete(ctx, pid, accessClaims.Sub.ID)

	if err != nil {
		return resp.InternalServerError("error updating project").WithTracePrefix("database-error").AddTrace(err)
	}

	return nil
}
