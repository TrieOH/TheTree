package handlers

import (
	"net/http"
	"payssage/models"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
)

// @Router /providers/{provider}/connect [post]
func (h *Handlers) Connect(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	provider, err := req.Path("provider").StringRequired()
	if fun.Bail(w, err) {
		return
	}
	var payload models.ConnectRequest
	if bind.BailInto(w, req, &payload) {
		return
	}
	if !payload.Flow.IsValid() {
		fun.BadRequest("invalid flow, either collector or seller").Send(w)
		return
	}
	if payload.Flow == models.OAuthFlowSeller && payload.WalletID == nil {
		fun.BadRequest("flow seller requires wallet id").Send(w)
		return
	}
	if payload.Flow == models.OAuthFlowCollector && payload.WalletID != nil {
		fun.BadRequest("flow collector cannot have wallet id").Send(w)
		return
	}
	redirectURL, err := h.commands.Connect(r.Context(), payload.ToInput(provider))
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, redirectURL)
}
