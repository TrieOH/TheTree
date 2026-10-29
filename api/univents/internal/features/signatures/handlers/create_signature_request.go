package handlers

import (
	"context"

	"univents/internal/openapi"
	"univents/models"
)

func (h *Handlers) CreateSignatureRequest(ctx context.Context, req openapi.CreateSignatureRequestRequestObject) (openapi.CreateSignatureRequestResponseObject, error) {
	expiresInDays := 7 // default when omitted
	if req.Body.ExpiresInDays != nil {
		expiresInDays = *req.Body.ExpiresInDays
	}
	request, err := h.ops.CreateRequest(ctx, models.CreateSignatureRequestInput{
		EditionID:       req.EditionId,
		SignatoryName:   req.Body.SignatoryName,
		SignatoryTitle:  req.Body.SignatoryTitle,
		SignatoryEmail:  &req.Body.SignatoryEmail,
		SignatoryUserID: req.Body.SignatoryUserId,
		ExpiresInDays:   expiresInDays,
	})
	if err != nil {
		return nil, err
	}
	return openapi.CreateSignatureRequest201JSONResponse(*request), nil
}
