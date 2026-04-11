package forms

import (
	"TrieForms/internal/shared/validation"
	"net/http"

	_ "TrieForms/internal/shared/contracts"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
)

type Handler struct {
	commands *CommandService
	queries  *QueryService
}

func NewFormsHandler(
	commands *CommandService,
	queries *QueryService,
) *Handler {
	return &Handler{
		commands: commands,
		queries:  queries,
	}
}

type CreateFormRequest struct {
	Title string `json:"title" validate:"required"`
}

// Create godoc
// @Summary Create a form
// @Description Creates a form in the given project.
// @Tags forms
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param project_id path string true "Project ID"
// @Param request body CreateFormRequest true "Form title"
// @Success 201 {object} contracts.Form "Form created successfully"
// @Failure 400 {object} contracts.ErrorResponse
// @Failure 401 {object} contracts.ErrorResponse
// @Failure 404 {object} contracts.ErrorResponse
// @Failure 500 {object} contracts.ErrorResponse
// @Router /projects/{project_id}/forms [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	projectID, rs := validation.GetUUID(r, "project_id")
	if rs == nil {
		rs.Send(w)
		return
	}

	var req CreateFormRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		resp.Error(err).Send(w)
		return
	}

	form, err := h.commands.Create(r.Context(), req.Title, projectID)
	if err != nil {
		resp.Error(err).Send(w)
		return
	}

	resp.Created().WithData(form).Send(w)
}

// List godoc
// @Summary List Forms
// @Description Lists all Forms for the given project
// @Tags forms
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param project_id path string true "Project ID"
// @Success 200 {array} contracts.Form "Forms retrieved successfully"
// @Failure 401 {object} contracts.ErrorResponse
// @Failure 404 {object} contracts.ErrorResponse
// @Failure 500 {object} contracts.ErrorResponse
// @Router /projects/{project_id}/forms [get]
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	projectID, rs := validation.GetUUID(r, "project_id")
	if rs == nil {
		rs.Send(w)
		return
	}

	forms, err := h.queries.List(r.Context(), projectID)
	if err != nil {
		resp.Error(err).Send(w)
		return
	}

	resp.OK().WithData(forms).Send(w)
}
