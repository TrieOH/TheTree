package editions

import (
	"net/http"
	"univents/internal/core/application/activity/commands"
	"univents/internal/core/domain"
	"univents/internal/core/interfaces/http/dto"
	"univents/internal/shared/validation"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
)

type Handler struct {
	commands *commands.CommandService
	//queries  *queries.QueryService
}

func NewActivitiesHandler(
	commands *commands.CommandService,
	// queries *queries.QueryService,
) *Handler {
	return &Handler{
		commands: commands,
		//queries:  queries,
	}
}

// Create godoc
// @Summary Create a new activity
// @Description Creates a new activity for an edition.
// @Tags activities
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param event_id path string true "Event ID"
// @Param edition_id path string true "Edition ID"
// @Param request body dto.CreateActivityRequest true "Activity creation request"
// @Success 201 {object} object "Activity created successfully"
// @Failure 400 {object} swag.ErrorResponse
// @Failure 401 {object} swag.ErrorResponse
// @Failure 404 {object} swag.ErrorResponse
// @Failure 500 {object} swag.ErrorResponse
// @Router /events/{event_id}/editions/{edition_id}/activities [post]
func (handler *Handler) Create(w http.ResponseWriter, r *http.Request) {
	editionID, rs := validation.GetUUID(r, "edition_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	var req dto.CreateActivityRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		resp.FromError(err).Send(w)
		return
	}

	in := domain.CreateActivitySpec{
		EditionScopeID: req.EditionScopeID,
		EditionID:      editionID,
		Title:          req.Title,
		Description:    req.Description,
		Location:       req.Location,
		StartsAt:       req.StartsAt,
		EndsAt:         req.EndsAt,
		PresenterName:  req.PresenterName,
		TokenCost:      req.TokenCost,
		HasCapacity:    req.HasCapacity,
		Capacity:       req.Capacity,
		Difficulty:     req.Difficulty,
	}

	ctx := r.Context()
	out, err := handler.commands.Create(ctx, in)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.Created().WithData(out).Send(w)
}
