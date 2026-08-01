// Package actors implements the StrictServerInterface methods for the
// actors feature. Auth and validation run in the strict middleware stack;
// these handlers are pure domain logic + envelope.
package actors

import (
	"context"
	"time"

	"IdentityX/internal/openapi"
	"IdentityX/internal/services"
	"IdentityX/models"
)

const module = "IdentityX"

type Handlers struct {
	ops *services.Actors
}

func New(ops *services.Actors) *Handlers { return &Handlers{ops: ops} }

func (h *Handlers) ListActors(ctx context.Context, req openapi.ListActorsRequestObject) (openapi.ListActorsResponseObject, error) {
	actors, err := h.ops.List(ctx, req.ProjectId)
	if err != nil {
		return nil, err
	}
	return openapi.ListActors200JSONResponse{
		Code: 200, Data: &actors, Timestamp: time.Now(), Module: module,
	}, nil
}

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
	return openapi.CreateActor201JSONResponse{
		Code: 201, Data: actor, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) GetActor(ctx context.Context, req openapi.GetActorRequestObject) (openapi.GetActorResponseObject, error) {
	actor, err := h.ops.GetByID(ctx, req.ActorId, req.ProjectId)
	if err != nil {
		return nil, err
	}
	return openapi.GetActor200JSONResponse{
		Code: 200, Data: actor, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) GetActorByEmail(ctx context.Context, req openapi.GetActorByEmailRequestObject) (openapi.GetActorByEmailResponseObject, error) {
	actor, err := h.ops.GetByEmail(ctx, req.ActorEmail, req.ProjectId)
	if err != nil {
		return nil, err
	}
	return openapi.GetActorByEmail200JSONResponse{
		Code: 200, Data: actor, Timestamp: time.Now(), Module: module,
	}, nil
}
