package handlers

import (
	"net/http"
	"univents/contracts"
	"univents/internal/shared/validation"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

type CreateEventRequest struct {
	OrganizationID *uuid.UUID `json:"organization_id"`
	Name           string     `json:"name" validate:"required,min=2"`
	Acronym        *string    `json:"acronym"`
	Slug           string     `json:"slug" validate:"required,min=2"`
	Tagline        *string    `json:"tagline"`
	Description    *string    `json:"description"`
	IsSeries       bool       `json:"is_series"`
	LogoUrl        *string    `json:"logo_url"`
	BannerUrl      *string    `json:"banner_url"`
	ContactEmail   *string    `json:"contact_email" validate:"required,email"`
}

func (handler *Handlers) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateEventRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		fun.Error(err).Send(w)
		return
	}

	in := contracts.CreateEventSpec{
		OrganizationID: req.OrganizationID,
		Name:           req.Name,
		Acronym:        req.Acronym,
		Slug:           req.Slug,
		Tagline:        req.Tagline,
		Description:    req.Description,
		IsSeries:       req.IsSeries,
		LogoUrl:        req.LogoUrl,
		ContactEmail:   req.ContactEmail,
	}

	ctx := r.Context()
	out, err := handler.commands.CreateEvent(ctx, in)
	if err != nil {
		fun.Error(err).Send(w)
		return
	}

	fun.Created().WithData(out).Send(w)
}
