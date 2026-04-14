package checkpoints

import (
	"net/http"
	"time"
	"univents/internal/shared/contracts"
	"univents/internal/shared/validation"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
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

type CreateCheckpointRequest struct {
	StartsAt   *time.Time                 `json:"starts_at"`
	EndsAt     *time.Time                 `json:"ends_at"`
	Name       string                     `json:"name"`
	Type       contracts.CheckpointType   `json:"type"`
	AccessMode contracts.CheckpointAccess `json:"access_mode" validate:"required"`
}

// Create godoc
// @Summary Create a new checkpoint
// @Description Creates a new checkpoint for an edition.
// @Tags checkpoints
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param event_id path string true "Event ID"
// @Param edition_id path string true "Edition ID"
// @Param request body CreateCheckpointRequest true "Checkpoint creation request"
// @Success 201 {object} object "Checkpoint created successfully"
// @Failure 400 {object} contracts.ErrorResponse
// @Failure 401 {object} contracts.ErrorResponse
// @Failure 404 {object} contracts.ErrorResponse
// @Failure 500 {object} contracts.ErrorResponse
// @Router /events/{event_id}/editions/{edition_id}/checkpoints [post]
func (handler *Handler) Create(w http.ResponseWriter, r *http.Request) {
	editionID, rs := validation.GetUUID(r, "edition_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	var req CreateCheckpointRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		resp.FromError(err).Send(w)
		return
	}

	in := contracts.CreateCheckpointSpec{
		EditionID:  editionID,
		StartsAt:   req.StartsAt,
		EndsAt:     req.EndsAt,
		Name:       req.Name,
		Type:       req.Type,
		AccessMode: req.AccessMode,
	}

	ctx := r.Context()
	out, err := handler.commands.Create(ctx, in)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.Created().WithData(out).Send(w)
}

// List godoc
// @Summary List all edition checkpoints
// @Description if user has permission checkpoints:read list all edition checkpoints
// @Tags checkpoints
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
// @Router /events/{event_id}/editions/{edition_id}/checkpoints [get]
func (handler *Handler) List(w http.ResponseWriter, r *http.Request) {
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
