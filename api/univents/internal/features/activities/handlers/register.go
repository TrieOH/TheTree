package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
)

func (handler *Handlers) Register(w http.ResponseWriter, r *http.Request) {
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
