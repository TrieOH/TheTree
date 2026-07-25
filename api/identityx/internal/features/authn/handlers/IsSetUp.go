package handlers

import (
	"lib/globals"
	"net/http"

	"github.com/MintzyG/fun"
)

func (h *Handlers) IsSetup(w http.ResponseWriter, _ *http.Request) {
	if globals.SetupComplete() {
		fun.ServiceUnavailable("setup already complete").Send(w)
		return
	}
	fun.NoContent().Send(w)
}
