package handlers

import (
	"context"

	"univents/internal/features/signatures/service"
	"univents/internal/openapi"
	"univents/models"
)

type Handlers struct {
	ops *service.Operations
}

func New(ops *service.Operations) *Handlers { return &Handlers{ops: ops} }

func (h *Handlers) CancelSignatureRequest(ctx context.Context, req openapi.CancelSignatureRequestRequestObject) (openapi.CancelSignatureRequestResponseObject, error) {
	err := h.ops.CancelRequest(ctx, req.RequestId, req.Body.Reason)
	if err != nil {
		return nil, err
	}
	return openapi.CancelSignatureRequest204Response{}, nil
}

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
	return openapi.CreateSignature201JSONResponse(*signature), nil
}

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

func (h *Handlers) DeleteSignature(ctx context.Context, req openapi.DeleteSignatureRequestObject) (openapi.DeleteSignatureResponseObject, error) {
	err := h.ops.Delete(ctx, req.SignatureId)
	if err != nil {
		return nil, err
	}
	return openapi.DeleteSignature204Response{}, nil
}

func (h *Handlers) DenySignatureRequest(ctx context.Context, req openapi.DenySignatureRequestRequestObject) (openapi.DenySignatureRequestResponseObject, error) {
	err := h.ops.DenyRequest(ctx, req.Params.Token, req.Body.Reason)
	if err != nil {
		return nil, err
	}
	return openapi.DenySignatureRequest204Response{}, nil
}

func (h *Handlers) FulfillSignatureRequest(ctx context.Context, req openapi.FulfillSignatureRequestRequestObject) (openapi.FulfillSignatureRequestResponseObject, error) {
	signature, err := h.ops.FulfillRequest(ctx, req.Params.Token, req.Body.ImageUrl)
	if err != nil {
		return nil, err
	}
	return openapi.FulfillSignatureRequest201JSONResponse(*signature), nil
}

func (h *Handlers) GetSignature(ctx context.Context, req openapi.GetSignatureRequestObject) (openapi.GetSignatureResponseObject, error) {
	signature, err := h.ops.GetByID(ctx, req.SignatureId)
	if err != nil {
		return nil, err
	}
	return openapi.GetSignature200JSONResponse(*signature), nil
}

func (h *Handlers) GetSignatureRequest(ctx context.Context, req openapi.GetSignatureRequestRequestObject) (openapi.GetSignatureRequestResponseObject, error) {
	request, err := h.ops.GetRequestByID(ctx, req.RequestId)
	if err != nil {
		return nil, err
	}
	return openapi.GetSignatureRequest200JSONResponse(*request), nil
}

func (h *Handlers) ListEditionSignatureRequests(ctx context.Context, req openapi.ListEditionSignatureRequestsRequestObject) (openapi.ListEditionSignatureRequestsResponseObject, error) {
	requests, err := h.ops.ListRequestsByEdition(ctx, req.EditionId)
	if err != nil {
		return nil, err
	}
	return openapi.ListEditionSignatureRequests200JSONResponse{Items: requests}, nil
}

func (h *Handlers) ListEditionSignatures(ctx context.Context, req openapi.ListEditionSignaturesRequestObject) (openapi.ListEditionSignaturesResponseObject, error) {
	signatures, err := h.ops.ListByEdition(ctx, req.EditionId)
	if err != nil {
		return nil, err
	}
	return openapi.ListEditionSignatures200JSONResponse{Items: signatures}, nil
}

func (h *Handlers) RevokeSignature(ctx context.Context, req openapi.RevokeSignatureRequestObject) (openapi.RevokeSignatureResponseObject, error) {
	err := h.ops.RevokeSignature(ctx, req.Params.Token)
	if err != nil {
		return nil, err
	}
	return openapi.RevokeSignature204Response{}, nil
}
