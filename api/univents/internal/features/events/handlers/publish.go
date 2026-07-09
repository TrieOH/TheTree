package handlers

import (
	"net/http"
	"univents/internal/shared/validation"

	"github.com/MintzyG/fun"
)

func (handler *Handlers) Publish(w http.ResponseWriter, r *http.Request) {
	eventID, rs := validation.GetUUID(r, "event_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	ctx := r.Context()
	err := handler.commands.PublishEvent(ctx, eventID)
	if err != nil {
		fun.Error(err).Send(w)
		return
	}

	fun.OK().Send(w)
}
