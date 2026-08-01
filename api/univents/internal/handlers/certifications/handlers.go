// Package certifications implements the StrictServerInterface methods for
// the certifications feature.
package certifications

import (
	"context"
	"encoding/json"
	"time"

	idx "sdk/identityx"
	"univents/internal/openapi"
	"univents/internal/services"
	"univents/models"

	"github.com/MintzyG/fun"
)

const module = "Univents"

type Handlers struct {
	ops *services.Certs
}

func New(ops *services.Certs) *Handlers { return &Handlers{ops: ops} }

func (h *Handlers) ListCertificationTemplates(ctx context.Context, req openapi.ListCertificationTemplatesRequestObject) (openapi.ListCertificationTemplatesResponseObject, error) {
	templates, err := h.ops.ListTemplates(ctx, req.EditionId)
	if err != nil {
		return nil, err
	}
	return openapi.ListCertificationTemplates200JSONResponse{
		Code: 200, Data: &templates, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) GetCertificationTemplate(ctx context.Context, req openapi.GetCertificationTemplateRequestObject) (openapi.GetCertificationTemplateResponseObject, error) {
	template, err := h.ops.GetTemplateByID(ctx, req.TemplateId)
	if err != nil {
		return nil, err
	}
	return openapi.GetCertificationTemplate200JSONResponse{
		Code: 200, Data: template, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) CreateCertificationTemplate(ctx context.Context, req openapi.CreateCertificationTemplateRequestObject) (openapi.CreateCertificationTemplateResponseObject, error) {
	design, err := json.Marshal(req.Body.DesignData)
	if err != nil {
		return nil, fun.ErrBadRequest("invalid design_data")
	}
	template, err := h.ops.CreateTemplate(ctx, models.CreateCertificationTemplateInput{
		EditionID:   req.EditionId,
		Kind:        req.Body.Kind,
		Name:        req.Body.Name,
		Description: req.Body.Description,
		DesignData:  design,
	})
	if err != nil {
		return nil, err
	}
	return openapi.CreateCertificationTemplate201JSONResponse{
		Code: 201, Data: template, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) UpdateCertificationTemplate(ctx context.Context, req openapi.UpdateCertificationTemplateRequestObject) (openapi.UpdateCertificationTemplateResponseObject, error) {
	design, err := json.Marshal(req.Body.DesignData)
	if err != nil {
		return nil, fun.ErrBadRequest("invalid design_data")
	}
	template, err := h.ops.UpdateTemplate(ctx, models.UpdateCertificationTemplateInput{
		TemplateID:  req.TemplateId,
		Kind:        req.Body.Kind,
		Name:        req.Body.Name,
		Description: req.Body.Description,
		DesignData:  design,
	})
	if err != nil {
		return nil, err
	}
	return openapi.UpdateCertificationTemplate200JSONResponse{
		Code: 200, Data: template, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) DeleteCertificationTemplate(ctx context.Context, req openapi.DeleteCertificationTemplateRequestObject) (openapi.DeleteCertificationTemplateResponseObject, error) {
	err := h.ops.DeleteTemplate(ctx, req.TemplateId)
	if err != nil {
		return nil, err
	}
	return openapi.DeleteCertificationTemplate204Response{}, nil
}

func (h *Handlers) LinkCertificationTemplate(ctx context.Context, req openapi.LinkCertificationTemplateRequestObject) (openapi.LinkCertificationTemplateResponseObject, error) {
	err := h.ops.LinkCertTemplate(ctx, req.TemplateId, req.Body.ProgramId)
	if err != nil {
		return nil, err
	}
	return openapi.LinkCertificationTemplate201Response{}, nil
}

func (h *Handlers) UnlinkCertificationTemplate(ctx context.Context, req openapi.UnlinkCertificationTemplateRequestObject) (openapi.UnlinkCertificationTemplateResponseObject, error) {
	err := h.ops.UnlinkCertTemplate(ctx, req.TemplateId, req.Body.ProgramId)
	if err != nil {
		return nil, err
	}
	return openapi.UnlinkCertificationTemplate204Response{}, nil
}

func (h *Handlers) ListCertificationTemplateLinks(ctx context.Context, req openapi.ListCertificationTemplateLinksRequestObject) (openapi.ListCertificationTemplateLinksResponseObject, error) {
	links, err := h.ops.ListCertTemplateLinks(ctx, req.TemplateId)
	if err != nil {
		return nil, err
	}
	return openapi.ListCertificationTemplateLinks200JSONResponse{
		Code: 200, Data: &links, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) VerifyCertification(ctx context.Context, req openapi.VerifyCertificationRequestObject) (openapi.VerifyCertificationResponseObject, error) {
	cert, err := h.ops.GetCertByHash(ctx, req.Hash)
	if err != nil {
		return nil, err
	}
	resp := models.VerifyCertResponse{
		Valid:      cert.Valid,
		TemplateID: cert.TemplateID,
		Cert:       cert,
	}
	return openapi.VerifyCertification200JSONResponse{
		Code: 200, Data: &resp, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) GetCertification(ctx context.Context, req openapi.GetCertificationRequestObject) (openapi.GetCertificationResponseObject, error) {
	cert, err := h.ops.GetCertByID(ctx, req.CertId)
	if err != nil {
		return nil, err
	}
	return openapi.GetCertification200JSONResponse{
		Code: 200, Data: cert, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) ListEditionCertifications(ctx context.Context, req openapi.ListEditionCertificationsRequestObject) (openapi.ListEditionCertificationsResponseObject, error) {
	certs, err := h.ops.ListCertsByEdition(ctx, req.EditionId)
	if err != nil {
		return nil, err
	}
	return openapi.ListEditionCertifications200JSONResponse{
		Code: 200, Data: &certs, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) ListMyCertifications(ctx context.Context, _ openapi.ListMyCertificationsRequestObject) (openapi.ListMyCertificationsResponseObject, error) {
	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}
	certs, err := h.ops.ListCertsByUser(ctx, ident.Sub.ID)
	if err != nil {
		return nil, err
	}
	return openapi.ListMyCertifications200JSONResponse{
		Code: 200, Data: &certs, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) InvalidateCertification(ctx context.Context, req openapi.InvalidateCertificationRequestObject) (openapi.InvalidateCertificationResponseObject, error) {
	err := h.ops.InvalidateCert(ctx, req.CertId, &req.Body.Reason)
	if err != nil {
		return nil, err
	}
	return openapi.InvalidateCertification204Response{}, nil
}

func (h *Handlers) ListCertificationEmissionErrors(ctx context.Context, req openapi.ListCertificationEmissionErrorsRequestObject) (openapi.ListCertificationEmissionErrorsResponseObject, error) {
	errors, err := h.ops.ListEmissionErrors(ctx, req.EditionId)
	if err != nil {
		return nil, err
	}
	return openapi.ListCertificationEmissionErrors200JSONResponse{
		Code: 200, Data: &errors, Timestamp: time.Now(), Module: module,
	}, nil
}
