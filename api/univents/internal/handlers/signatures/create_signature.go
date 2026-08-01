package signatures

import (
	"context"
	"time"

	"univents/internal/openapi"
	"univents/models"
)

func (h *Handlers) CreateSignature(ctx context.Context, req openapi.CreateSignatureRequestObject) (openapi.CreateSignatureResponseObject, error) {
	signature, err := h.ops.Create(ctx, models.AddSignatureInput{
		EditionID:       req.EditionId,
		SignatoryName:   req.Body.SignatoryName,
		SignatoryTitle:  req.Body.SignatoryTitle,
		SignatoryEmail:  req.Body.SignatoryEmail,
		SignatoryUserID: req.Body.SignatoryUserId,
		ImageURL:        req.Body.ImageUrl,
	})
	if err != nil {
		return nil, err
	}
	return openapi.CreateSignature201JSONResponse{
		Code: 201, Data: signature, Timestamp: time.Now(), Module: module,
	}, nil
}
