package activities

import (
	"net/http"

	"univents/internal/shared/contracts"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
	"github.com/go-chi/chi/v5"
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
	r.Route("/events/{event_id}/editions/{edition_id}/activities", func(r chi.Router) {
		r.Get("/", h.List)
		r.Use(jwt)
		r.Post("/", h.Create)
		r.Get("/admin", h.ListAdmin)
		r.Route("/{event_id}", func(r chi.Router) {
			r.Post("/publish", h.Publish)
			r.Post("/register", h.Register)
			r.Post("/unregister", h.Unregister)
			r.Get("/records", h.ListRecords)
			r.Post("/records/{record_id}", h.MarkAttendance)
		})
	})
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
	req := fun.From(r)
	editionID, err := req.Path("edition_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	var payload contracts.CreateActivityRequest
	if bind.BailInto(w, req, &payload) {
		return
	}
	activity, err := handler.commands.Create(r.Context(), payload.ToSpec(editionID))
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, activity, http.StatusCreated)
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
	req := fun.From(r)
	activityID, err := req.Path("activity_id").UUID()
	err = handler.commands.Publish(r.Context(), activityID)
	if fun.Bail(w, err) {
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
	req := fun.From(r)
	editionID, err := req.Path("edition_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	activities, err := handler.queries.List(r.Context(), editionID)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, activities)
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
	req := fun.From(r)
	editionID, err := req.Path("edition_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	activities, err := handler.queries.AdminList(r.Context(), editionID)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, activities)
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
	req := fun.From(r)
	activityID, err := req.Path("activity_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	err = handler.commands.Register(r.Context(), activityID)
	if fun.Bail(w, err) {
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
	req := fun.From(r)
	activityID, err := req.Path("activity_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	err = handler.commands.Unregister(r.Context(), activityID)
	if fun.Bail(w, err) {
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
	req := fun.From(r)
	activityID, err := req.Path("activity_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	recordID, err := req.Path("record_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	err = handler.commands.MarkAttendance(r.Context(), activityID, recordID)
	if fun.Bail(w, err) {
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
	req := fun.From(r)
	activityID, err := req.Path("activity_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	records, err := handler.commands.ListRecords(r.Context(), activityID)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, records)
}
