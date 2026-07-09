package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
)

func (handler *Handlers) Unregister(w http.ResponseWriter, r *http.Request) {
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
