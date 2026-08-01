package certifications

import (
	"encoding/json"

	"context"
	"time"

	"github.com/MintzyG/fun"

	"univents/internal/openapi"
	"univents/models"
)

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
