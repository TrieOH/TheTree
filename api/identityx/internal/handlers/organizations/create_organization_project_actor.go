package organizations

import (
	"context"
	"time"

	"IdentityX/internal/openapi"
	"IdentityX/models"
)

func (h *Handlers) CreateOrganizationProjectActor(ctx context.Context, req openapi.CreateOrganizationProjectActorRequestObject) (openapi.CreateOrganizationProjectActorResponseObject, error) {
	actor, err := h.ops.CreateProjectActor(ctx, models.CreateActorInput{
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
