package tos

import (
	"context"
	"time"

	"IdentityX/internal/openapi"
	"IdentityX/models"
)

func (h *Handlers) GetProjectTos(ctx context.Context, req openapi.GetProjectTosRequestObject) (openapi.GetProjectTosResponseObject, error) {
	tos, err := h.ops.Get(ctx, req.ProjectId)
	if err != nil {
		return nil, err
	}
	data := toOpenTos(*tos)
	return openapi.GetProjectTos200JSONResponse{
		Code: 200, Data: &data, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) PostProjectTos(ctx context.Context, req openapi.PostProjectTosRequestObject) (openapi.PostProjectTosResponseObject, error) {
	tos, err := h.ops.Create(ctx, req.ProjectId, req.Body.Content)
	if err != nil {
		return nil, err
	}
	data := toOpenTos(*tos)
	return openapi.PostProjectTos200JSONResponse{
		Code: 200, Data: &data, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) PutProjectTos(ctx context.Context, req openapi.PutProjectTosRequestObject) (openapi.PutProjectTosResponseObject, error) {
	tos, err := h.ops.Update(ctx, req.ProjectId, req.Body.Content)
	if err != nil {
		return nil, err
	}
	data := toOpenTos(*tos)
	return openapi.PutProjectTos200JSONResponse{
		Code: 200, Data: &data, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) PostTosAcceptance(ctx context.Context, req openapi.PostTosAcceptanceRequestObject) (openapi.PostTosAcceptanceResponseObject, error) {
	acceptance, err := h.ops.Accept(ctx, req.ProjectId)
	if err != nil {
		return nil, err
	}
	data := toOpenAcceptance(models.TosAcceptanceWithActor{TosAcceptance: *acceptance})
	return openapi.PostTosAcceptance200JSONResponse{
		Code: 200, Data: &data, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) ListTosAcceptances(ctx context.Context, req openapi.ListTosAcceptancesRequestObject) (openapi.ListTosAcceptancesResponseObject, error) {
	acceptances, err := h.ops.ListAcceptances(ctx, req.ProjectId)
	if err != nil {
		return nil, err
	}
	data := make([]openapi.TosAcceptance, 0, len(acceptances))
	for _, a := range acceptances {
		data = append(data, toOpenAcceptance(a))
	}
	return openapi.ListTosAcceptances200JSONResponse{
		Code: 200, Data: &data, Timestamp: time.Now(), Module: module,
	}, nil
}
