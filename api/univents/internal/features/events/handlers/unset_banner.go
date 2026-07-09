package handlers

import (
	"net/http"
	"univents/internal/shared/validation"

	"github.com/MintzyG/fun"
)

func (handler *Handlers) UnsetBanner(w http.ResponseWriter, r *http.Request) {
	eventID, rs := validation.GetUUID(r, "event_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	ctx := r.Context()
	product, err := handler.commands.UnsetBanner(ctx, eventID)
	if err != nil {
		fun.Error(err).Send(w)
		return
	}

	fun.OK("Logo unset").WithData(product).Send(w)
}
