package handlers

import (
	"net/http"
	"univents/internal/shared/validation"

	"github.com/MintzyG/fun"
)

func (handler *Handlers) DisconnectPaymentAccount(w http.ResponseWriter, r *http.Request) {
	_, rs := validation.GetUUID(r, "event_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	editionID, rs := validation.GetUUID(r, "edition_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	ctx := r.Context()
	err := handler.commands.DisconnectPayments(ctx, editionID)
	if err != nil {
		fun.Error(err).Send(w)
		return
	}

	fun.OK().Send(w)
}
