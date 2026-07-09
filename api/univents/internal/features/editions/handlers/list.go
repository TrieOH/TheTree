package handlers

import (
	"net/http"
	"univents/internal/shared/validation"

	"github.com/MintzyG/fun"
)

func (handler *Handlers) List(w http.ResponseWriter, r *http.Request) {
	eventID, rs := validation.GetUUID(r, "event_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	ctx := r.Context()
	out, err := handler.queries.List(ctx, eventID)
	if err != nil {
		fun.Error(err).Send(w)
		return
	}

	fun.Created().WithData(out).Send(w)
}
