package handlers

import (
	"net/http"
	"univents/models"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
)

func (handler *Handlers) CreateVariant(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	productID, err := req.Path("product_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	var payload models.CreateProductVariantRequest
	if bind.BailInto(w, req, &payload) {
		return
	}
	variant, err := handler.commands.CreateVariant(r.Context(), payload.ToInput(productID))
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, variant, http.StatusCreated)
}
