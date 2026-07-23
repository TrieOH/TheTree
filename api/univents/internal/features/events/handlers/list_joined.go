package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
)

func (handler *Handlers) ListJoined(w http.ResponseWriter, r *http.Request) {
	events, err := handler.queries.ListJoined(r.Context())
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, events)
}
