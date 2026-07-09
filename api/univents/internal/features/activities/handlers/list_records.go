package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
)

func (handler *Handlers) ListRecords(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	activityID, err := req.Path("activity_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	records, err := handler.queries.ListRecords(r.Context(), activityID)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, records)
}
