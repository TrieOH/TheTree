package actors

import (
	"context"
	"time"

	"IdentityX/internal/openapi"
	"IdentityX/models"
)

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
