package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
)

func (handler *Handlers) GetRequestByID(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	id, err := req.Path("request_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	request, err := handler.queries.GetRequestByID(r.Context(), id)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, request)
}
