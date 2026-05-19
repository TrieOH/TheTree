package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
)

// ListNamespaces godoc
// @Summary Lists the namespaces you have access too
// @Description Lists the namespaces you have access to, this includes your own and the ones you were invited too
// @Tags namespaces
// @ID namespaces_listnamespaces
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} fun.Response{data=[]models.Namespace}
// @Failure 400 {object} fun.Response
// @Failure 401 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Router /namespaces [get]
func (h *Handler) ListNamespaces(w http.ResponseWriter, r *http.Request) {
	namespaces, err := h.queries.ListNamespaces(r.Context())
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, namespaces)
}
