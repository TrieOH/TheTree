package tickets

import (
	"net/http"
	"univents/internal/commerce/application/ticket/commands"
	"univents/internal/commerce/domain"
	"univents/internal/commerce/interfaces/http/dtos"
	"univents/internal/shared/validation"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
)

type TicketsHandler struct {
	commands *commands.CommandService
	//queries  *queries.QueryService
}

func NewTicketsHandler(
	commands *commands.CommandService,
	// queries *queries.QueryService,
) *TicketsHandler {
	return &TicketsHandler{
		commands: commands,
		//queries:  queries,
	}
}

// Create godoc
// @Summary Create a new ticket
// @Description Creates a new ticket for an edition.
// @Tags tickets
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param event_id path string true "Event ID"
// @Param edition_id path string true "Edition ID"
// @Param request body dtos.CreateTicketRequest true "Ticket creation request"
// @Success 201 {object} object "Ticket created successfully"
// @Failure 400 {object} swag.ErrorResponse
// @Failure 401 {object} swag.ErrorResponse
// @Failure 404 {object} swag.ErrorResponse
// @Failure 500 {object} swag.ErrorResponse
// @Router /events/{event_id}/editions/{edition_id}/tickets [post]
func (handler *TicketsHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dtos.CreateTicketRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		resp.FromError(err).Send(w)
		return
	}

	editionID, rs := validation.GetUUID(r, "edition_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	in := domain.CreateTicketSpec{
		EditionScopeID: req.EditionScopeID,
		EditionID:      editionID,
		Name:           req.Name,
		Description:    req.Description,
	}

	ctx := r.Context()
	out, err := handler.commands.Create(ctx, in)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.Created().WithData(out).Send(w)
}
