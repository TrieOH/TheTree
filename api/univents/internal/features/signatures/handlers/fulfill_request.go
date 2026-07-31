package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
)

type fulfillRequestPayload struct {
	ImageURL string `json:"image_url" validate:"required,url"`
}

func (handler *Handlers) FulfillRequest(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	token, err := req.Query("token").StringRequired()
	if fun.Bail(w, err) {
		return
	}
	var payload fulfillRequestPayload
	if bind.BailInto(w, req, &payload) {
		return
	}
	signature, err := handler.commands.FulfillRequest(r.Context(), token, payload.ImageURL)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, signature, http.StatusCreated)
}
