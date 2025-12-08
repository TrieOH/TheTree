package service

import (
	"GoAuth/internal/models"
	"GoAuth/internal/repository"
	"GoAuth/internal/utils"
	"context"
	"log"
	"net/http"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

func (s *AuthService) CreateProject(r *http.Request, project models.Project) (*models.Project, *resp.Response) {
	accessClaims, err := models.GetAccessClaims(r)
	if err != nil {
		return nil, resp.InternalServerError().AddTrace(err)
	}

	project.OwnerId = accessClaims.Sub.ID
	project.IsActive = true

	pubKey, privKey, err := utils.GenerateEd25519Keys()
	if err != nil {
		return nil, resp.InternalServerError("error generating public and private project keys").AddTrace(err)
	}

	project.PubKey = pubKey
	project.PrivKey = []byte(privKey)
	log.Println(privKey)
	log.Println(pubKey)

	dbProject, err := s.queries.CreateProject(r.Context(), repository.CreateProjectParams{
		ProjectName: project.ProjectName,
		OwnerID:     project.OwnerId,
		Metadata:    project.Metadata,
		IsActive:    project.IsActive,
		PubKey:      project.PubKey,
		PrivKey:     string(project.PrivKey),
	})

	if err != nil {
		return nil, resp.InternalServerError("error creating project").WithTracePrefix("database-error").AddTrace(err)
	}

	var projectDTO models.Project
	copier.Copy(&projectDTO, &dbProject)

	return &projectDTO, nil
}

func (s *AuthService) GetProjectByID(ctx context.Context, r *http.Request, projectID string) (*models.Project, *resp.Response) {
	accessClaims, err := models.GetAccessClaims(r)
	if err != nil {
		return nil, resp.InternalServerError().AddTrace(err)
	}

	pid, err := uuid.Parse(projectID)
	if err != nil {
		return nil, resp.InternalServerError("error parsing project uuid").AddTrace(err)
	}

	dbProject, err := s.queries.GetProjectById(ctx, repository.GetProjectByIdParams{
		OwnerID: accessClaims.Sub.ID,
		ID:      pid,
	})

	if err != nil {
		return nil, resp.InternalServerError("error fetching project").WithTracePrefix("database-error").AddTrace(err)
	}

	var projectDTO models.Project
	copier.Copy(&projectDTO, &dbProject)

	return &projectDTO, nil
}

func (s *AuthService) ListProjects(ctx context.Context, r *http.Request) ([]models.Project, *resp.Response) {
	accessClaims, err := models.GetAccessClaims(r)
	if err != nil {
		return nil, resp.InternalServerError().AddTrace(err)
	}

	dbProjects, err := s.queries.ListProjects(ctx, accessClaims.Sub.ID)

	if err != nil {
		return nil, resp.InternalServerError("error fetching projects").WithTracePrefix("database-error").AddTrace(err)
	}

	var projectDTOs []models.Project
	copier.Copy(&projectDTOs, &dbProjects)

	return projectDTOs, nil
}

func (s *AuthService) GetProjectKeysByID(ctx context.Context, r *http.Request, projectId string) (*models.ProjectKeys, *resp.Response) {
	accessClaims, err := models.GetAccessClaims(r)
	if err != nil {
		return nil, resp.InternalServerError().AddTrace(err)
	}

	pid, err := uuid.Parse(projectId)
	if err != nil {
		return nil, resp.BadRequest("error parsing project id").AddTrace(err)
	}

	dbProject, err := s.queries.GetProjectKeysById(ctx, repository.GetProjectKeysByIdParams{
		OwnerID: accessClaims.Sub.ID,
		ID:      pid,
	})

	if err != nil {
		return nil, resp.InternalServerError("error fetching project keys").WithTracePrefix("database-error").AddTrace(err)
	}

	var keys models.ProjectKeys
	copier.Copy(&keys, &dbProject)

	return &keys, nil
}

func (s *AuthService) GetProjectJWKS(ctx context.Context, projectId string) (map[string]any, *resp.Response) {
	pid, err := uuid.Parse(projectId)
	if err != nil {
		return nil, resp.BadRequest("error parsing project id").AddTrace(err)
	}

	pubKey, err := s.queries.GetProjectPublicKeyById(ctx, pid)
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

func (s *AuthService) UpdateProjectByID(ctx context.Context, r *http.Request, projectId string, project models.Project) (*models.Project, *resp.Response) {
	accessClaims, err := models.GetAccessClaims(r)
	if err != nil {
		return nil, resp.InternalServerError().AddTrace(err)
	}

	newProject, rs := s.GetProjectByID(ctx, r, projectId)
	if rs != nil {
		return nil, rs
	}

	if project.ProjectName != "" {
		newProject.ProjectName = project.ProjectName
	}
	newProject.Metadata = project.Metadata

	updatedProject, err := s.queries.UpdateProject(ctx, repository.UpdateProjectParams{
		OwnerID:     accessClaims.Sub.ID,
		ID:          newProject.ID,
		ProjectName: newProject.ProjectName,
		Metadata:    newProject.Metadata,
	})

	if err != nil {
		return nil, resp.InternalServerError("error updating project").WithTracePrefix("database-error").AddTrace(err)
	}

	var projectDTO models.Project
	copier.Copy(&projectDTO, &updatedProject)

	return &projectDTO, nil
}

func (s *AuthService) DeleteProjectByID(ctx context.Context, r *http.Request, projectId string) *resp.Response {
	accessClaims, err := models.GetAccessClaims(r)
	if err != nil {
		return resp.InternalServerError().AddTrace(err)
	}

	pid, err := uuid.Parse(projectId)
	if err != nil {
		return resp.BadRequest("error parsing project id").AddTrace(err)
	}

	err = s.queries.DeleteProject(ctx, repository.DeleteProjectParams{
		ID:      pid,
		OwnerID: accessClaims.Sub.ID,
	})

	if err != nil {
		return resp.InternalServerError("error updating project").WithTracePrefix("database-error").AddTrace(err)
	}

	return nil
}
