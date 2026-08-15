package organizations

import (
	"context"
	"time"

	"IdentityX/internal/openapi"
	"IdentityX/models"
)

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
