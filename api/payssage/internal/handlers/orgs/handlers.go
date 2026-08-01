// Package orgs implements the StrictServerInterface methods for the orgs
// feature.
package orgs

import (
	"context"
	"time"

	"payssage/internal/openapi"
	"payssage/internal/services"
	"payssage/models"
)

const module = "Payssage"

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

func (h *Handlers) GetOrganizationMemberByID(ctx context.Context, req openapi.GetOrganizationMemberByIDRequestObject) (openapi.GetOrganizationMemberByIDResponseObject, error) {
	member, err := h.ops.GetMemberByID(ctx, req.MemberId, req.OrganizationId)
	if err != nil {
		return nil, err
	}
	return openapi.GetOrganizationMemberByID200JSONResponse{
		Code: 200, Data: member, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) GetOrganizationMemberByEmail(ctx context.Context, req openapi.GetOrganizationMemberByEmailRequestObject) (openapi.GetOrganizationMemberByEmailResponseObject, error) {
	member, err := h.ops.GetMemberByEmail(ctx, req.MemberEmail, req.OrganizationId)
	if err != nil {
		return nil, err
	}
	return openapi.GetOrganizationMemberByEmail200JSONResponse{
		Code: 200, Data: member, Timestamp: time.Now(), Module: module,
	}, nil
}
