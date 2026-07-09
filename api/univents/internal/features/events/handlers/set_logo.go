package handlers

import (
	"net/http"
	"univents/internal/shared/validation"

	"github.com/MintzyG/fun"
)

func (handler *Handlers) SetLogo(w http.ResponseWriter, r *http.Request) {
	eventID, rs := validation.GetUUID(r, "event_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	var req ImageURLRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		fun.Error(err).Send(w)
		return
	}

	ctx := r.Context()
	product, err := handler.commands.SetLogo(ctx, eventID, req.URL)
	if err != nil {
		fun.Error(err).Send(w)
		return
	}

	fun.OK("Logo set").WithData(product).Send(w)
}
