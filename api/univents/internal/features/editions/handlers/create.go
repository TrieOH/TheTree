package handlers

import (
	"net/http"
	"univents/contracts"
	"univents/internal/shared/validation"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
)

func (handler *Handlers) Create(w http.ResponseWriter, r *http.Request) {
	eventID, rs := validation.GetUUID(r, "event_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	var req contracts.CreateEditionRequest
	if bind.BailInto(w, fun.From(r), &req) {
		return
	}

	in := contracts.CreateEditionSpec{
		EventID:              eventID,
		Type:                 req.Type,
		EditionName:          req.EditionName,
		Tagline:              req.Tagline,
		Description:          req.Description,
		RegistrationOpensAt:  req.RegistrationOpensAt,
		RegistrationClosesAt: req.RegistrationClosesAt,
		StartsAt:             req.StartsAt,
		EndsAt:               req.EndsAt,
		Timezone:             req.Timezone,
		LocationName:         req.LocationName,
		LocationAddress:      req.LocationAddress,
		LogoUrl:              req.LogoUrl,
		BannerUrl:            req.BannerUrl,
		ContactEmail:         req.ContactEmail,
		ContactPhone:         req.ContactPhone,
		OrganizerName:        req.OrganizerName,
	}

	ctx := r.Context()
	out, err := handler.commands.Create(ctx, in)
	if err != nil {
		fun.Error(err).Send(w)
		return
	}

	fun.Created().WithData(out).Send(w)
}
