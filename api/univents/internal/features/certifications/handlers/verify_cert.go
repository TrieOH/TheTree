package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
	"github.com/go-chi/chi/v5"
)

func (handler *Handlers) VerifyCert(w http.ResponseWriter, r *http.Request) {
	hash := chi.URLParam(r, "hash")
	cert, err := handler.queries.GetCertByHash(r.Context(), hash)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, map[string]interface{}{
		"valid":       cert.Valid,
		"template_id": cert.TemplateID,
		"cert":        cert,
	})
}
