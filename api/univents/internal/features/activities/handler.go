package activities

import (
	"net/http"

	"univents/contracts"

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
		r.With(jwt).Post("/", h.Create)
		r.With(jwt).Get("/admin", h.ListAdmin)
		r.With(jwt).Route("/{event_id}", func(r chi.Router) {
			r.Post("/publish", h.Publish)
			r.Post("/register", h.Register)
			r.Post("/unregister", h.Unregister)
			r.Get("/records", h.ListRecords)
			r.Post("/records/{record_id}", h.MarkAttendance)
		})
	})
}

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

func (handler *Handler) Publish(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	activityID, err := req.Path("activity_id").UUID()
	err = handler.commands.Publish(r.Context(), activityID)
	if fun.Bail(w, err) {
		return
	}
	fun.OK().Send(w)
}

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
