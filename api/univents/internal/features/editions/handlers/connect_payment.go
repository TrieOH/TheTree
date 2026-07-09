package handlers

import (
	"net/http"
	"univents/internal/shared/validation"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

func (handler *Handlers) ConnectPaymentAccount(w http.ResponseWriter, r *http.Request) {
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

	triePaymentsCredentialID := r.URL.Query().Get("credential_id")
	if triePaymentsCredentialID == "" {
		fun.BadRequest("missing credential_id").Send(w)
		return
	}

	credID, err := uuid.Parse(triePaymentsCredentialID)
	if err != nil {
		fun.BadRequest("invalid credential_id: " + err.Error()).Send(w)
		return
	}

	provider := r.URL.Query().Get("provider")
	if provider == "" {
		fun.BadRequest("missing provider").Send(w)
		return
	}

	publicKey := r.URL.Query().Get("public_key")
	if publicKey == "" {
		fun.BadRequest("missing public_key").Send(w)
		return
	}

	ctx := r.Context()
	err = handler.commands.ConnectPayments(ctx, credID, editionID, provider, publicKey)
	if err != nil {
		fun.Error(err).Send(w)
		return
	}

	fun.OK("Payment account connected successfully").Send(w)
}
