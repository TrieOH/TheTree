package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
	"github.com/google/uuid"
)

func (handler *Handlers) LinkCertTemplate(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	templateID, err := req.Path("template_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	var payload struct {
		ProgramID string `json:"program_id" validate:"required,uuid"`
	}
	if bind.BailInto(w, req, &payload) {
		return
	}
	err = handler.commands.LinkCertTemplate(r.Context(), templateID, uuid.MustParse(payload.ProgramID))
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, nil, http.StatusCreated)
}
