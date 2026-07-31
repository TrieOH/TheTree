package handlers

import (
	"io"
	"net/http"
	"payssage/models"

	"github.com/MintzyG/fun"
	"github.com/go-chi/chi/v5"
)

func (h *Handlers) Receive(w http.ResponseWriter, r *http.Request) {
	providerName := chi.URLParam(r, "provider")

	rawBody, err := io.ReadAll(r.Body)
	if err != nil {
		fun.BadRequest("failed to read webhook body").Send(w)
		return
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(r.Body)

	err = h.commands.Receive(r.Context(), models.ReceiveWebhookInput{
		Provider: providerName,
		Request:  r,
		RawBody:  rawBody,
	})
	if fun.Bail(w, err) {
		return
	}

	// Always 200 once accepted, MP (and most providers) treat non-2xx
	// as "retry me," which you don't want once you've already processed
	// or intentionally ignored an event.
	fun.OK().Send(w)
}
