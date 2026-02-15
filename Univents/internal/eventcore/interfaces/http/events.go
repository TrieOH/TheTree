package http

import (
	"net/http"
	"univents/internal/eventcore/application/commands"
	"univents/internal/eventcore/application/queries"
	"univents/internal/eventcore/domain"
	"univents/internal/eventcore/interfaces/http/dto"
	"univents/internal/shared/validation"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
)

type EventsHandler struct {
	commands *commands.CommandService
	queries  *queries.QueryService
}

func NewEventsHandler(
	commands *commands.CommandService,
	queries *queries.QueryService,
) *EventsHandler {
	return &EventsHandler{
		commands: commands,
		queries:  queries,
	}
}

func (handler *EventsHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateEventRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		resp.FromError(err).Send(w)
		return
	}

	in := domain.Event{
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
	out, err := handler.commands.Create(ctx, in)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.Created().WithData(out).Send(w)
}

func (handler *EventsHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	out, err := handler.queries.List(ctx)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK().WithData(out).Send(w)
}

func (handler *EventsHandler) Publish(w http.ResponseWriter, r *http.Request) {
	eventID, rs := validation.GetUUID(r, "event_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	ctx := r.Context()
	err := handler.commands.Publish(ctx, eventID)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK().Send(w)
}
