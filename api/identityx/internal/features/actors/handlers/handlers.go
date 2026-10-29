package handlers

import (
	"context"

	"IdentityX/internal/features/actors/service"
	"IdentityX/internal/openapi"
	"IdentityX/models"
)

type Handlers struct {
	ops *service.Operations
}

func New(ops *service.Operations) *Handlers { return &Handlers{ops: ops} }

func (h *Handlers) CreateActor(ctx context.Context, req openapi.CreateActorRequestObject) (openapi.CreateActorResponseObject, error) {
	actor, err := h.ops.Create(ctx, models.CreateActorInput{
		ProjectID:  &req.ProjectId,
		AuthMethod: req.Body.AuthMethod,
		Type:       req.Body.Type,
		Email:      req.Body.Email,
	})
	if err != nil {
		return nil, err
	}
	return openapi.CreateActor201JSONResponse(*actor), nil
}

func (h *Handlers) GetActor(ctx context.Context, req openapi.GetActorRequestObject) (openapi.GetActorResponseObject, error) {
	actor, err := h.ops.GetByID(ctx, req.ActorId, req.ProjectId)
	if err != nil {
		return nil, err
	}
	return openapi.GetActor200JSONResponse(*actor), nil
}

func (h *Handlers) GetActorByEmail(ctx context.Context, req openapi.GetActorByEmailRequestObject) (openapi.GetActorByEmailResponseObject, error) {
	actor, err := h.ops.GetByEmail(ctx, req.ActorEmail, req.ProjectId)
	if err != nil {
		return nil, err
	}
	return openapi.GetActorByEmail200JSONResponse(*actor), nil
}

func (h *Handlers) ListActors(ctx context.Context, req openapi.ListActorsRequestObject) (openapi.ListActorsResponseObject, error) {
	actors, err := h.ops.List(ctx, req.ProjectId)
	if err != nil {
		return nil, err
	}
	return openapi.ListActors200JSONResponse{Items: actors}, nil
}
