// Package organizations implements the StrictServerInterface methods for
// the organizations feature.
package organizations

import (
	"context"
	"time"

	"IdentityX/internal/openapi"
	"IdentityX/internal/services"
	"IdentityX/models"
)

const module = "IdentityX"

type Handlers struct {
	ops *services.Organizations
}

func New(ops *services.Organizations) *Handlers { return &Handlers{ops: ops} }

func (h *Handlers) ListOrganizations(ctx context.Context, _ openapi.ListOrganizationsRequestObject) (openapi.ListOrganizationsResponseObject, error) {
	orgs, err := h.ops.ListOrgs(ctx)
	if err != nil {
		return nil, err
	}
	return openapi.ListOrganizations200JSONResponse{
		Code: 200, Data: &orgs, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) CreateOrganization(ctx context.Context, req openapi.CreateOrganizationRequestObject) (openapi.CreateOrganizationResponseObject, error) {
	org, err := h.ops.Create(ctx, models.CreateOrganizationInput{
		Name: req.Body.Name,
		Slug: req.Body.Slug,
	})
	if err != nil {
		return nil, err
	}
	return openapi.CreateOrganization201JSONResponse{
		Code: 201, Data: org, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) ListOrganizationMembers(ctx context.Context, req openapi.ListOrganizationMembersRequestObject) (openapi.ListOrganizationMembersResponseObject, error) {
	members, err := h.ops.ListMembers(ctx, req.OrganizationId)
	if err != nil {
		return nil, err
	}
	return openapi.ListOrganizationMembers200JSONResponse{
		Code: 200, Data: &members, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) AddOrganizationMember(ctx context.Context, req openapi.AddOrganizationMemberRequestObject) (openapi.AddOrganizationMemberResponseObject, error) {
	err := h.ops.AddMember(ctx, models.AddOrganizationMemberInput{
		ActorEmail:     req.Body.ActorEmail,
		Role:           req.Body.Role,
		OrganizationID: req.OrganizationId,
	})
	if err != nil {
		return nil, err
	}
	return openapi.AddOrganizationMember201Response{}, nil
}

func (h *Handlers) RemoveOrganizationMember(ctx context.Context, req openapi.RemoveOrganizationMemberRequestObject) (openapi.RemoveOrganizationMemberResponseObject, error) {
	err := h.ops.RemoveMember(ctx, models.RemoveOrganizationMemberInput{
		ActorEmail:     req.Body.ActorEmail,
		OrganizationID: req.OrganizationId,
	})
	if err != nil {
		return nil, err
	}
	return openapi.RemoveOrganizationMember200JSONResponse{
		Code: 200, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) ListOrganizationProjects(ctx context.Context, req openapi.ListOrganizationProjectsRequestObject) (openapi.ListOrganizationProjectsResponseObject, error) {
	projects, err := h.ops.ListOrgProjects(ctx, req.OrganizationId)
	if err != nil {
		return nil, err
	}
	return openapi.ListOrganizationProjects200JSONResponse{
		Code: 200, Data: &projects, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) CreateOrganizationProject(ctx context.Context, req openapi.CreateOrganizationProjectRequestObject) (openapi.CreateOrganizationProjectResponseObject, error) {
	project, err := h.ops.CreateProject(ctx, models.CreateOrgProjectInput{
		OrganizationID: req.OrganizationId,
		Name:           req.Body.Name,
		Domain:         req.Body.Domain,
		BrandSlug:      req.Body.BrandSlug,
	})
	if err != nil {
		return nil, err
	}
	return openapi.CreateOrganizationProject201JSONResponse{
		Code: 201, Data: project, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) ListOrganizationProjectActors(ctx context.Context, req openapi.ListOrganizationProjectActorsRequestObject) (openapi.ListOrganizationProjectActorsResponseObject, error) {
	actors, err := h.ops.ListProjectActors(ctx, req.OrgId, req.ProjectId)
	if err != nil {
		return nil, err
	}
	return openapi.ListOrganizationProjectActors200JSONResponse{
		Code: 200, Data: &actors, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) CreateOrganizationProjectActor(ctx context.Context, req openapi.CreateOrganizationProjectActorRequestObject) (openapi.CreateOrganizationProjectActorResponseObject, error) {
	actor, err := h.ops.CreateProjectActor(ctx, req.OrgId, models.CreateActorInput{
		ProjectID:  &req.ProjectId,
		AuthMethod: req.Body.AuthMethod,
		Type:       req.Body.Type,
		Email:      req.Body.Email,
	})
	if err != nil {
		return nil, err
	}
	return openapi.CreateOrganizationProjectActor201JSONResponse{
		Code: 201, Data: actor, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) ListOrganizationProjectMembers(ctx context.Context, req openapi.ListOrganizationProjectMembersRequestObject) (openapi.ListOrganizationProjectMembersResponseObject, error) {
	members, err := h.ops.ListOrgProjectMembers(ctx, req.OrganizationId, req.ProjectId)
	if err != nil {
		return nil, err
	}
	return openapi.ListOrganizationProjectMembers200JSONResponse{
		Code: 200, Data: &members, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) AddOrganizationProjectMember(ctx context.Context, req openapi.AddOrganizationProjectMemberRequestObject) (openapi.AddOrganizationProjectMemberResponseObject, error) {
	err := h.ops.AddProjectMember(ctx, models.AddOrgProjectMemberInput{
		ActorEmail:     req.Body.ActorEmail,
		Role:           req.Body.Role,
		OrganizationID: req.OrganizationId,
		ProjectID:      req.ProjectId,
	})
	if err != nil {
		return nil, err
	}
	return openapi.AddOrganizationProjectMember201Response{}, nil
}

func (h *Handlers) RemoveOrganizationProjectMember(ctx context.Context, req openapi.RemoveOrganizationProjectMemberRequestObject) (openapi.RemoveOrganizationProjectMemberResponseObject, error) {
	err := h.ops.RemoveProjectMember(ctx, models.RemoveOrgProjectMemberInput{
		ActorEmail:     req.Body.ActorEmail,
		OrganizationID: req.OrganizationId,
		ProjectID:      req.ProjectId,
	})
	if err != nil {
		return nil, err
	}
	return openapi.RemoveOrganizationProjectMember200JSONResponse{
		Code: 200, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) GetOrganizationProjectActor(ctx context.Context, req openapi.GetOrganizationProjectActorRequestObject) (openapi.GetOrganizationProjectActorResponseObject, error) {
	actor, err := h.ops.GetActorByID(ctx, req.ActorId, req.OrganizationId, req.ProjectId)
	if err != nil {
		return nil, err
	}
	return openapi.GetOrganizationProjectActor200JSONResponse{
		Code: 200, Data: actor, Timestamp: time.Now(), Module: module,
	}, nil
}
