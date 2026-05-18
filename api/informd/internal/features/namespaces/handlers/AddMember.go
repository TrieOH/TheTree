package handlers

import (
	"Informd/models"
	"net/http"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
)

// AddMember godoc
// @Summary Add a namespace member
// @Description Lets you add a member to the namespace as a viewer, editor or admin
// @Tags namespaces
// @ID namespaces_addmember
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.AddNamespaceMemberRequest true "Member details"
// @Success 201 {object} fun.Response
// @Failure 400 {object} fun.Response
// @Failure 401 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Router /namespaces/{namespace_id}/members [post]
func (h *Handler) AddMember(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	namespaceID, err := req.Path("namespace_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	var payload models.AddNamespaceMemberRequest
	if bind.BailInto(w, req, &payload) {
		return
	}
	err = h.commands.AddMember(r.Context(), payload.ToInput(namespaceID))
	if fun.Bail(w, err) {
		return
	}
	fun.Created().Send(w)
}
