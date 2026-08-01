// Package projects implements the StrictServerInterface methods for the
// projects feature.
package projects

import (
	"context"
	"time"

	"IdentityX/internal/openapi"
	"IdentityX/internal/services"
	"IdentityX/models"
)

const module = "IdentityX"

type Handlers struct {
	ops *services.Projects
}

func New(ops *services.Projects) *Handlers { return &Handlers{ops: ops} }

func (h *Handlers) ListProjects(ctx context.Context, _ openapi.ListProjectsRequestObject) (openapi.ListProjectsResponseObject, error) {
	projects, err := h.ops.ListProjects(ctx)
	if err != nil {
		return nil, err
	}
	return openapi.ListProjects200JSONResponse{
		Code: 200, Data: &projects, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) CreateProject(ctx context.Context, req openapi.CreateProjectRequestObject) (openapi.CreateProjectResponseObject, error) {
	project, err := h.ops.Create(ctx, models.CreateProjectInput{
		Name:      req.Body.Name,
		Domain:    req.Body.Domain,
		BrandSlug: req.Body.BrandSlug,
	})
	if err != nil {
		return nil, err
	}
	return openapi.CreateProject201JSONResponse{
		Code: 201, Data: project, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) ListProjectMembers(ctx context.Context, req openapi.ListProjectMembersRequestObject) (openapi.ListProjectMembersResponseObject, error) {
	members, err := h.ops.ListMembers(ctx, req.ProjectId)
	if err != nil {
		return nil, err
	}
	return openapi.ListProjectMembers200JSONResponse{
		Code: 200, Data: &members, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) AddProjectMember(ctx context.Context, req openapi.AddProjectMemberRequestObject) (openapi.AddProjectMemberResponseObject, error) {
	err := h.ops.AddMember(ctx, models.AddProjectMemberInput{
		ActorEmail: req.Body.ActorEmail,
		Role:       req.Body.Role,
		ProjectID:  req.ProjectId,
	})
	if err != nil {
		return nil, err
	}
	return openapi.AddProjectMember201Response{}, nil
}

func (h *Handlers) RemoveProjectMember(ctx context.Context, req openapi.RemoveProjectMemberRequestObject) (openapi.RemoveProjectMemberResponseObject, error) {
	err := h.ops.RemoveMember(ctx, models.RemoveProjectMemberInput{
		ActorEmail: req.Body.ActorEmail,
		ProjectID:  req.ProjectId,
	})
	if err != nil {
		return nil, err
	}
	return openapi.RemoveProjectMember200JSONResponse{
		Code: 200, Timestamp: time.Now(), Module: module,
	}, nil
}
