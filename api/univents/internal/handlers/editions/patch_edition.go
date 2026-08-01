package editions

import (
	"context"
	"time"

	"univents/internal/openapi"
	"univents/models"
)

func (h *Handlers) PatchEdition(ctx context.Context, req openapi.PatchEditionRequestObject) (openapi.PatchEditionResponseObject, error) {
	edition, err := h.ops.Patch(ctx, models.PatchEditionInput{
		EditionID:           req.EditionId,
		Name:                req.Body.Name,
		Slug:                req.Body.Slug,
		Tagline:             req.Body.Tagline,
		Description:         req.Body.Description,
		RegistrationOpensAt: req.Body.RegistrationOpensAt,
		StartsAt:            req.Body.StartsAt,
		EndsAt:              req.Body.EndsAt,
		LocationName:        req.Body.LocationName,
		LocationAddress:     req.Body.LocationDescription,
		LogoURL:             req.Body.LogoUrl,
		BannerURL:           req.Body.BannerUrl,
		ContactEmail:        req.Body.ContactEmail,
	})
	if err != nil {
		return nil, err
	}
	return openapi.PatchEdition200JSONResponse{
		Code: 200, Data: edition, Timestamp: time.Now(), Module: module,
	}, nil
}
