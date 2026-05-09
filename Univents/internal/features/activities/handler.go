package activities

import (
	"net/http"
	"time"
	"univents/internal/shared/contracts"
	"univents/internal/shared/validation"

	"github.com/MintzyG/fun"
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

type CreateActivityRequest struct {
	Title         string                     `json:"title" validate:"required,min=3"`
	Description   *string                    `json:"description"`
	Location      string                     `json:"location"`
	StartsAt      time.Time                  `json:"starts_at" validate:"required"`
	EndsAt        time.Time                  `json:"ends_at" validate:"required"`
	PresenterName *string                    `json:"presenter_name"`
	TokenCost     int                        `json:"token_cost" validate:"gte=0"`
	HasCapacity   bool                       `json:"has_capacity"`
	Capacity      int                        `json:"capacity" validate:"gte=0"`
	Difficulty    *contracts.DifficultyLevel `json:"difficulty"`
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
// @Param request body CreateActivityRequest true "Activity creation request"
// @Success 201 {object} object "Activity created successfully"
// @Failure 400 {object} contracts.ErrorResponse
// @Failure 401 {object} contracts.ErrorResponse
// @Failure 404 {object} contracts.ErrorResponse
// @Failure 500 {object} contracts.ErrorResponse
// @Router /events/{event_id}/editions/{edition_id}/activities [post]
func (handler *Handler) Create(w http.ResponseWriter, r *http.Request) {
	editionID, rs := validation.GetUUID(r, "edition_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	var req CreateActivityRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		fun.Error(err).Send(w)
		return
	}

	in := contracts.CreateActivitySpec{
		EditionID:     editionID,
		Title:         req.Title,
		Description:   req.Description,
		Location:      req.Location,
		StartsAt:      req.StartsAt,
		EndsAt:        req.EndsAt,
		PresenterName: req.PresenterName,
		TokenCost:     req.TokenCost,
		HasCapacity:   req.HasCapacity,
		Capacity:      req.Capacity,
		Difficulty:    req.Difficulty,
	}

	ctx := r.Context()
	out, err := handler.commands.Create(ctx, in)
	if err != nil {
		fun.Error(err).Send(w)
		return
	}

	fun.Created().WithData(out).Send(w)
}

// Publish godoc
// @Summary publishes an activity
// @Description Publishes an activity making it publicly available.
// @Tags activities
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param event_id path string true "Event ID"
// @Param edition_id path string true "Edition ID"
// @Param activity_id path string true "Activity ID"
// @Success 200 {object} object "Activity published successfully"
// @Failure 400 {object} contracts.ErrorResponse
// @Failure 401 {object} contracts.ErrorResponse
// @Failure 404 {object} contracts.ErrorResponse
// @Failure 500 {object} contracts.ErrorResponse
// @Router /events/{event_id}/editions/{edition_id}/activities/{activity_id}/publish [post]
func (handler *Handler) Publish(w http.ResponseWriter, r *http.Request) {
	activityID, rs := validation.GetUUID(r, "activity_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	ctx := r.Context()
	err := handler.commands.Publish(ctx, activityID)
	if err != nil {
		fun.Error(err).Send(w)
		return
	}

	fun.OK().Send(w)
}

// List godoc
// @Summary List all edition activities
// @Description List all publicly available activities of the edition
// @Tags activities
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
// @Router /events/{event_id}/editions/{edition_id}/activities [get]
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

// ListAdmin godoc
// @Summary List all edition activities
// @Description if user has permission activities:read list all edition activities
// @Tags activities
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
// @Router /events/{event_id}/editions/{edition_id}/activities/admin [get]
func (handler *Handler) ListAdmin(w http.ResponseWriter, r *http.Request) {
	editionID, rs := validation.GetUUID(r, "edition_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	ctx := r.Context()
	out, err := handler.queries.AdminList(ctx, editionID)
	if err != nil {
		fun.Error(err).Send(w)
		return
	}

	fun.OK().WithData(out).Send(w)
}

// Register godoc
// @Summary register the user to an activity
// @Description Registers the user to the specified activity if they have activities:attend permission on it
// @Tags activities
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param event_id path string true "Event ID"
// @Param edition_id path string true "Edition ID"
// @Param activity_id path string true "Activity ID"
// @Success 200 {object} object "Registered Successfully"
// @Failure 400 {object} contracts.ErrorResponse
// @Failure 401 {object} contracts.ErrorResponse
// @Failure 404 {object} contracts.ErrorResponse
// @Failure 500 {object} contracts.ErrorResponse
// @Router /events/{event_id}/editions/{edition_id}/activities/{activity_id}/register [post]
func (handler *Handler) Register(w http.ResponseWriter, r *http.Request) {
	activityID, rs := validation.GetUUID(r, "activity_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	ctx := r.Context()
	err := handler.commands.Register(ctx, activityID)
	if err != nil {
		fun.Error(err).Send(w)
		return
	}

	fun.OK("Registered Successfully").Send(w)
}

// Unregister godoc
// @Summary unregisters the user from an activity
// @Description Unregisters the user from the specified activity if they are registered on it
// @Tags activities
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param event_id path string true "Event ID"
// @Param edition_id path string true "Edition ID"
// @Param activity_id path string true "Activity ID"
// @Success 200 {object} object "Unregistered Successfully"
// @Failure 400 {object} contracts.ErrorResponse
// @Failure 401 {object} contracts.ErrorResponse
// @Failure 404 {object} contracts.ErrorResponse
// @Failure 500 {object} contracts.ErrorResponse
// @Router /events/{event_id}/editions/{edition_id}/activities/{activity_id}/unregister [post]
func (handler *Handler) Unregister(w http.ResponseWriter, r *http.Request) {
	activityID, rs := validation.GetUUID(r, "activity_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	ctx := r.Context()
	err := handler.commands.Unregister(ctx, activityID)
	if err != nil {
		fun.Error(err).Send(w)
		return
	}

	fun.OK("Unregistered Successfully").Send(w)
}

// MarkAttendance godoc
// @Summary Marks attendance for a user in an activity
// @Description If you have attendance:mark on the activity and the record status is registered, marks it as completed
// @Tags activities
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param event_id path string true "Event ID"
// @Param edition_id path string true "Edition ID"
// @Param activity_id path string true "Activity ID"
// @Success 200 {object} object "Marked Attendance Successfully"
// @Failure 400 {object} contracts.ErrorResponse
// @Failure 401 {object} contracts.ErrorResponse
// @Failure 404 {object} contracts.ErrorResponse
// @Failure 500 {object} contracts.ErrorResponse
// @Router /events/{event_id}/editions/{edition_id}/activities/{activity_id}/records/{record_id} [post]
func (handler *Handler) MarkAttendance(w http.ResponseWriter, r *http.Request) {
	activityID, rs := validation.GetUUID(r, "activity_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	recordID, rs := validation.GetUUID(r, "record_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	ctx := r.Context()
	err := handler.commands.MarkAttendance(ctx, activityID, recordID)
	if err != nil {
		fun.Error(err).Send(w)
		return
	}

	fun.OK("Marked Attendance Successfully").Send(w)
}

// ListRecords godoc
// @Summary Lists attendance records of the activity
// @Description Lists attendance records of the activity if you have activities:manage on the activity
// @Tags activities
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param event_id path string true "Event ID"
// @Param edition_id path string true "Edition ID"
// @Param activity_id path string true "Activity ID"
// @Success 200 {object} object
// @Failure 400 {object} contracts.ErrorResponse
// @Failure 401 {object} contracts.ErrorResponse
// @Failure 404 {object} contracts.ErrorResponse
// @Failure 500 {object} contracts.ErrorResponse
// @Router /events/{event_id}/editions/{edition_id}/activities/{activity_id}/records [post]
func (handler *Handler) ListRecords(w http.ResponseWriter, r *http.Request) {
	activityID, rs := validation.GetUUID(r, "activity_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	ctx := r.Context()
	records, err := handler.commands.ListRecords(ctx, activityID)
	if err != nil {
		fun.Error(err).Send(w)
		return
	}

	fun.OK().WithData(records).Send(w)
}
