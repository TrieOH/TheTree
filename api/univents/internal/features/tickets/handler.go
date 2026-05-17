package tickets

import (
	"net/http"
	"univents/internal/shared/contracts"
	"univents/internal/shared/validation"

	"github.com/MintzyG/fun"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Handler struct {
	commands *CommandService
	queries  *QueryService
}

func NewHandler(
	commands *CommandService,
	queries *QueryService,
) *Handler {
	return &Handler{
		commands: commands,
		queries:  queries,
	}
}

func Routes(
	r *chi.Mux,
	h *Handler,
	jwt func(http.Handler) http.Handler,
) {
	r.Route("/events/{event_id}/editions/{edition_id}/tickets", func(r chi.Router) {
		r.Get("/", h.List)
		r.Use(jwt)
		r.Post("/", h.Create)
		r.Post("/{ticket_id}/permissions", h.AddPermission)
		r.Delete("/{ticket_id}/permissions/{permission_id}", h.RemovePermission)
	})
}

type CreateTicketRequest struct {
	Name        string  `json:"name" validate:"required,min=3"`
	Description *string `json:"description"`
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
// @Param request body CreateTicketRequest true "Ticket creation request"
// @Success 201 {object} object "Ticket created successfully"
// @Failure 400 {object} contracts.ErrorResponse
// @Failure 401 {object} contracts.ErrorResponse
// @Failure 404 {object} contracts.ErrorResponse
// @Failure 500 {object} contracts.ErrorResponse
// @Router /events/{event_id}/editions/{edition_id}/tickets [post]
func (handler *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateTicketRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		fun.Error(err).Send(w)
		return
	}

	editionID, rs := validation.GetUUID(r, "edition_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	in := contracts.CreateTicketSpec{
		EditionID:   editionID,
		Name:        req.Name,
		Description: req.Description,
	}

	ctx := r.Context()
	out, err := handler.commands.Create(ctx, in)
	if err != nil {
		fun.Error(err).Send(w)
		return
	}

	fun.Created().WithData(out).Send(w)
}

type AddTicketPermissionRequest struct {
	PermissionType contracts.PermissionType `json:"permission_type"`
	ActivityID     *uuid.UUID               `json:"activity_id"`
	ProductID      *uuid.UUID               `json:"product_id"`
	CheckpointID   *uuid.UUID               `json:"checkpoint_id"`
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
// @Param request body AddTicketPermissionRequest true "Ticket Permission attach request"
// @Success 201 {object} object "Ticket Permission added successfully"
// @Failure 400 {object} contracts.ErrorResponse
// @Failure 401 {object} contracts.ErrorResponse
// @Failure 404 {object} contracts.ErrorResponse
// @Failure 500 {object} contracts.ErrorResponse
// @Router /events/{event_id}/editions/{edition_id}/tickets/{ticket_id}/permissions [post]
func (handler *Handler) AddPermission(w http.ResponseWriter, r *http.Request) {
	var req AddTicketPermissionRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		fun.Error(err).Send(w)
		return
	}

	ticketID, rs := validation.GetUUID(r, "ticket_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	in := contracts.CreateTicketPermissionSpec{
		TicketID:       ticketID,
		PermissionType: req.PermissionType,
		ActivityID:     req.ActivityID,
		ProductID:      req.ProductID,
		CheckpointID:   req.CheckpointID,
	}

	ctx := r.Context()
	out, err := handler.commands.AddPermission(ctx, in)
	if err != nil {
		fun.Error(err).Send(w)
		return
	}

	fun.Created("Ticket Permission added successfully").WithData(out).Send(w)
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
// @Success 201 {object} object "Ticket Permission removed successfully"
// @Failure 400 {object} contracts.ErrorResponse
// @Failure 401 {object} contracts.ErrorResponse
// @Failure 404 {object} contracts.ErrorResponse
// @Failure 500 {object} contracts.ErrorResponse
// @Router /events/{event_id}/editions/{edition_id}/tickets/{ticket_id}/permissions/{permission_id} [delete]
func (handler *Handler) RemovePermission(w http.ResponseWriter, r *http.Request) {
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
	err := handler.commands.RemovePermission(ctx, permissionID, ticketID)
	if err != nil {
		fun.Error(err).Send(w)
		return
	}

	fun.OK("Ticket Permission removed successfully").Send(w)
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
// @Failure 400 {object} contracts.ErrorResponse
// @Failure 401 {object} contracts.ErrorResponse
// @Failure 404 {object} contracts.ErrorResponse
// @Failure 500 {object} contracts.ErrorResponse
// @Router /events/{event_id}/editions/{edition_id}/tickets [get]
func (handler *Handler) List(w http.ResponseWriter, r *http.Request) {
	editionID, rs := validation.GetUUID(r, "edition_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	ctx := r.Context()
	out, err := handler.queries.List(ctx, editionID)
	if err != nil {
		fun.Error(err).Send(w)
		return
	}

	fun.OK().WithData(out).Send(w)
}
