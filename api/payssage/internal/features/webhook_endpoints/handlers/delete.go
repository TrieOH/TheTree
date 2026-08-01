package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
)

func (h *Handlers) Delete(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	endpointID, err := req.Path("endpoint_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	err = h.ops.Delete(r.Context(), endpointID)
	if fun.Bail(w, err) {
		return
	}
	fun.OK().Send(w)
}
