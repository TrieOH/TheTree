// Package tos implements the project terms-of-service endpoints: the
// versioned document (get/post/put) and the consent ledger (accept/list).
package handlers

import (
	"context"

	"IdentityX/internal/features/tos/service"
	"IdentityX/internal/openapi"
	"IdentityX/models"
)

type Handlers struct {
	ops *service.Operations
}

func New(ops *service.Operations) *Handlers { return &Handlers{ops: ops} }

func toOpenTos(t models.TermsOfService) openapi.ProjectTos {
	return openapi.ProjectTos{
		Id:          t.ID,
		ProjectId:   t.ProjectID,
		Version:     t.Version,
		Content:     t.Content,
		EffectiveAt: t.EffectiveAt,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}

func toOpenAcceptance(a models.TosAcceptanceWithActor) openapi.TosAcceptance {
	return openapi.TosAcceptance{
		Id:          a.ID,
		ActorId:     a.ActorID,
		ProjectId:   a.ProjectID,
		TosVersion:  a.TosVersion,
		ContentHash: a.ContentHash,
		Source:      openapi.TosAcceptanceSource(a.Source),
		AcceptedAt:  a.AcceptedAt,
		ActorEmail:  a.ActorEmail,
	}
}

func (h *Handlers) GetProjectTos(ctx context.Context, req openapi.GetProjectTosRequestObject) (openapi.GetProjectTosResponseObject, error) {
	tos, err := h.ops.Get(ctx, req.ProjectId)
	if err != nil {
		return nil, err
	}
	data := toOpenTos(*tos)
	return openapi.GetProjectTos200JSONResponse(data), nil
}

func (h *Handlers) PostProjectTos(ctx context.Context, req openapi.PostProjectTosRequestObject) (openapi.PostProjectTosResponseObject, error) {
	tos, err := h.ops.Create(ctx, req.ProjectId, req.Body.Content)
	if err != nil {
		return nil, err
	}
	data := toOpenTos(*tos)
	return openapi.PostProjectTos200JSONResponse(data), nil
}

func (h *Handlers) PutProjectTos(ctx context.Context, req openapi.PutProjectTosRequestObject) (openapi.PutProjectTosResponseObject, error) {
	tos, err := h.ops.Update(ctx, req.ProjectId, req.Body.Content)
	if err != nil {
		return nil, err
	}
	data := toOpenTos(*tos)
	return openapi.PutProjectTos200JSONResponse(data), nil
}

func (h *Handlers) PostTosAcceptance(ctx context.Context, req openapi.PostTosAcceptanceRequestObject) (openapi.PostTosAcceptanceResponseObject, error) {
	acceptance, err := h.ops.Accept(ctx, req.ProjectId)
	if err != nil {
		return nil, err
	}
	data := toOpenAcceptance(models.TosAcceptanceWithActor{TosAcceptance: *acceptance})
	return openapi.PostTosAcceptance200JSONResponse(data), nil
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
	return openapi.ListTosAcceptances200JSONResponse{Items: data}, nil
}
