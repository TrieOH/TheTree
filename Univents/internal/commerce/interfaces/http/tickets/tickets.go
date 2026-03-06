package tickets

import (
	"net/http"
	"univents/internal/commerce/application/ticket/commands"
	"univents/internal/commerce/application/ticket/queries"
	"univents/internal/commerce/domain"
	"univents/internal/commerce/interfaces/http/dtos"
	"univents/internal/shared/validation"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
)

type TicketsHandler struct {
	commands *commands.CommandService
	queries  *queries.QueryService
}

func NewTicketsHandler(
	commands *commands.CommandService,
	queries *queries.QueryService,
) *TicketsHandler {
	return &TicketsHandler{
		commands: commands,
		queries:  queries,
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

// AddPermission godoc
// @Summary Adds a new permission to a ticket
// @Description Attaches a new permission to a ticket to allow user access to that object when they have the ticket
// @Tags tickets
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param event_id path string true "Event ID"
// @Param edition_id path string true "Edition ID"
// @Param ticket_id path string true "Ticket ID"
// @Param request body dtos.AddTicketPermissionRequest true "Ticket Permission attach request"
// @Success 201 {object} object "Ticket Permission added successfully"
// @Failure 400 {object} swag.ErrorResponse
// @Failure 401 {object} swag.ErrorResponse
// @Failure 404 {object} swag.ErrorResponse
// @Failure 500 {object} swag.ErrorResponse
// @Router /events/{event_id}/editions/{edition_id}/tickets/{ticket_id}/permissions [post]
func (handler *TicketsHandler) AddPermission(w http.ResponseWriter, r *http.Request) {
	var req dtos.AddTicketPermissionRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		resp.FromError(err).Send(w)
		return
	}

	ticketID, rs := validation.GetUUID(r, "ticket_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	in := domain.CreateTicketPermissionSpec{
		TicketScopeID:  req.TicketScopeID,
		TicketID:       ticketID,
		PermissionType: req.PermissionType,
		ActivityID:     req.ActivityID,
		ProductID:      req.ProductID,
		CheckpointID:   req.CheckpointID,
	}

	ctx := r.Context()
	out, err := handler.commands.AddPermission(ctx, in)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.Created("Ticket Permission added successfully").WithData(out).Send(w)
}

// RemovePermission godoc
// @Summary Removes a permission from a ticket
// @Description Removes a certain permission from a ticket therefore revoking access to all users that have that ticket from that permission
// @Tags tickets
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param event_id path string true "Event ID"
// @Param edition_id path string true "Edition ID"
// @Param ticket_id path string true "Ticket ID"
// @Param request body dtos.RemoveTicketPermissionRequest true "Ticket Permission removal request"
// @Success 201 {object} object "Ticket Permission removed successfully"
// @Failure 400 {object} swag.ErrorResponse
// @Failure 401 {object} swag.ErrorResponse
// @Failure 404 {object} swag.ErrorResponse
// @Failure 500 {object} swag.ErrorResponse
// @Router /events/{event_id}/editions/{edition_id}/tickets/{ticket_id}/permissions/{permission_id} [delete]
func (handler *TicketsHandler) RemovePermission(w http.ResponseWriter, r *http.Request) {
	var req dtos.RemoveTicketPermissionRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		resp.FromError(err).Send(w)
		return
	}

	ticketID, rs := validation.GetUUID(r, "ticket_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	permissionID, rs := validation.GetUUID(r, "permission_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	ctx := r.Context()
	err := handler.commands.RemovePermission(ctx, permissionID, ticketID, req.TicketScopeID)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK("Ticket Permission removed successfully").Send(w)
}

// List godoc
// @Summary List all edition tickets
// @Description if user has permission tickets:read list all edition tickets
// @Tags tickets
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param event_id path string true "Event ID"
// @Param edition_id path string true "Edition ID"
// @Success 201 {object} object
// @Failure 400 {object} swag.ErrorResponse
// @Failure 401 {object} swag.ErrorResponse
// @Failure 404 {object} swag.ErrorResponse
// @Failure 500 {object} swag.ErrorResponse
// @Router /events/{event_id}/editions/{edition_id}/tickets [get]
func (handler *TicketsHandler) List(w http.ResponseWriter, r *http.Request) {
	editionID, rs := validation.GetUUID(r, "edition_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	ctx := r.Context()
	out, err := handler.queries.List(ctx, editionID)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK().WithData(out).Send(w)
}
