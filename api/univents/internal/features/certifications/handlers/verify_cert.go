package handlers

import (
	"net/http"

	"univents/models"

	"github.com/MintzyG/fun"
	"github.com/go-chi/chi/v5"
)

func (handler *Handlers) VerifyCert(w http.ResponseWriter, r *http.Request) {
	hash := chi.URLParam(r, "hash")
	cert, err := handler.queries.GetCertByHash(r.Context(), hash)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, models.VerifyCertResponse{
		Valid:      cert.Valid,
		TemplateID: cert.TemplateID,
		Cert:       cert,
	})
}
