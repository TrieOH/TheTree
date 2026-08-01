package handlers

import (
	"net/http"
	"os"
	"payssage/models"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
)

func (h *Handlers) HardCreate(w http.ResponseWriter, r *http.Request) {
	if os.Getenv("TEST_MODE") != "true" {
		fun.ServiceUnavailable("test mode only").Send(w)
		return
	}

	var payload models.HardCreateIntentRequest
	if bind.BailInto(w, fun.From(r), &payload) {
		return
	}
	intent, err := h.ops.HardCreate(r.Context(), payload)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, intent, http.StatusCreated)
}
