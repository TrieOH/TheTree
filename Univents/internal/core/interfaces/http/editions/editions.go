package editions

import (
	"net/http"
	"univents/internal/core/application/edition/commands"
	"univents/internal/core/application/edition/queries"
	"univents/internal/core/domain"
	"univents/internal/core/interfaces/http/dto"
	"univents/internal/shared/validation"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
)

type Handler struct {
	commands *commands.CommandService
	queries  *queries.QueryService
}

func NewEditionsHandler(
	commands *commands.CommandService,
	queries *queries.QueryService,
) *Handler {
	return &Handler{
		commands: commands,
		queries:  queries,
	}
}

func (handler *Handler) Create(w http.ResponseWriter, r *http.Request) {
	eventID, rs := validation.GetUUID(r, "event_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	var req dto.CreateEditionRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		resp.FromError(err).Send(w)
		return
	}

	in := domain.CreateEditionSpec{
		EventID:              eventID,
		GoAuthEventScopeID:   req.GoAuthEventScopeID,
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
		resp.FromError(err).Send(w)
		return
	}

	resp.Created().WithData(out).Send(w)
}

func (handler *Handler) List(w http.ResponseWriter, r *http.Request) {
	eventID, rs := validation.GetUUID(r, "event_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	ctx := r.Context()
	out, err := handler.queries.ListEditions(ctx, eventID)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.Created().WithData(out).Send(w)
}

func (handler *Handler) Announce(w http.ResponseWriter, r *http.Request) {
	eventID, rs := validation.GetUUID(r, "event_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	editionID, rs := validation.GetUUID(r, "edition_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	ctx := r.Context()
	err := handler.commands.Announce(ctx, eventID, editionID)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.Created().Send(w)
}
