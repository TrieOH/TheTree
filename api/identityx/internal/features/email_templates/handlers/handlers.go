// Package email_templates implements the project-scoped email template
// override endpoints (list/get/put/delete).
package handlers

import (
	"context"

	"IdentityX/internal/features/email_templates/service"
	"IdentityX/internal/openapi"
	"IdentityX/models"
)

type Handlers struct {
	ops *service.Operations
}

func New(ops *service.Operations) *Handlers { return &Handlers{ops: ops} }

func toEffective(t models.EffectiveEmailTemplate) openapi.EffectiveEmailTemplate {
	return openapi.EffectiveEmailTemplate{
		Kind:    openapi.EmailTemplateKind(t.Kind),
		Subject: t.Subject,
		Body:    t.Body,
		Source:  openapi.EffectiveEmailTemplateSource(t.Source),
	}
}

func (h *Handlers) ListEmailTemplates(ctx context.Context, req openapi.ListEmailTemplatesRequestObject) (openapi.ListEmailTemplatesResponseObject, error) {
	templates, err := h.ops.List(ctx, req.ProjectId)
	if err != nil {
		return nil, err
	}
	data := make([]openapi.EffectiveEmailTemplate, 0, len(templates))
	for _, t := range templates {
		data = append(data, toEffective(t))
	}
	return openapi.ListEmailTemplates200JSONResponse{Items: data}, nil
}

func (h *Handlers) GetEmailTemplate(ctx context.Context, req openapi.GetEmailTemplateRequestObject) (openapi.GetEmailTemplateResponseObject, error) {
	template, err := h.ops.Get(ctx, req.ProjectId, models.EmailTemplateKind(req.Kind))
	if err != nil {
		return nil, err
	}
	effective := toEffective(*template)
	return openapi.GetEmailTemplate200JSONResponse(effective), nil
}

func (h *Handlers) PutEmailTemplate(ctx context.Context, req openapi.PutEmailTemplateRequestObject) (openapi.PutEmailTemplateResponseObject, error) {
	template, err := h.ops.Upsert(ctx, req.ProjectId, models.EmailTemplateKind(req.Kind), req.Body.Subject, req.Body.Body)
	if err != nil {
		return nil, err
	}
	effective := toEffective(*template)
	return openapi.PutEmailTemplate200JSONResponse(effective), nil
}

func (h *Handlers) DeleteEmailTemplate(ctx context.Context, req openapi.DeleteEmailTemplateRequestObject) (openapi.DeleteEmailTemplateResponseObject, error) {
	err := h.ops.Delete(ctx, req.ProjectId, models.EmailTemplateKind(req.Kind))
	if err != nil {
		return nil, err
	}
	return openapi.DeleteEmailTemplate204Response{}, nil
}
