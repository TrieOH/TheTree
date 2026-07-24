package handlers

import (
	"net/http"
	"univents/models"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
)

func (handler *Handlers) PatchVariant(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	variantID, err := req.Path("variant_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	var payload models.PatchProductVariantRequest
	if bind.BailInto(w, req, &payload) {
		return
	}
	variant, err := handler.commands.PatchVariant(r.Context(), payload.ToInput(variantID))
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, variant)
}
