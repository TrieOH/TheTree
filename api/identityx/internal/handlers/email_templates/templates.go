package email_templates

import (
	"context"
	"time"

	"IdentityX/internal/openapi"
	"IdentityX/models"
)

func (h *Handlers) ListEmailTemplates(ctx context.Context, req openapi.ListEmailTemplatesRequestObject) (openapi.ListEmailTemplatesResponseObject, error) {
	templates, err := h.ops.List(ctx, req.ProjectId)
	if err != nil {
		return nil, err
	}
	data := make([]openapi.EffectiveEmailTemplate, 0, len(templates))
	for _, t := range templates {
		data = append(data, toEffective(t))
	}
	return openapi.ListEmailTemplates200JSONResponse{
		Code: 200, Data: &data, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) GetEmailTemplate(ctx context.Context, req openapi.GetEmailTemplateRequestObject) (openapi.GetEmailTemplateResponseObject, error) {
	template, err := h.ops.Get(ctx, req.ProjectId, models.EmailTemplateKind(req.Kind))
	if err != nil {
		return nil, err
	}
	effective := toEffective(*template)
	return openapi.GetEmailTemplate200JSONResponse{
		Code: 200, Data: &effective, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) PutEmailTemplate(ctx context.Context, req openapi.PutEmailTemplateRequestObject) (openapi.PutEmailTemplateResponseObject, error) {
	template, err := h.ops.Upsert(ctx, req.ProjectId, models.EmailTemplateKind(req.Kind), req.Body.Subject, req.Body.Body)
	if err != nil {
		return nil, err
	}
	effective := toEffective(*template)
	return openapi.PutEmailTemplate200JSONResponse{
		Code: 200, Data: &effective, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) DeleteEmailTemplate(ctx context.Context, req openapi.DeleteEmailTemplateRequestObject) (openapi.DeleteEmailTemplateResponseObject, error) {
	err := h.ops.Delete(ctx, req.ProjectId, models.EmailTemplateKind(req.Kind))
	if err != nil {
		return nil, err
	}
	return openapi.DeleteEmailTemplate200JSONResponse{
		Code: 200, Timestamp: time.Now(), Module: module,
	}, nil
}
