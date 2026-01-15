package http

import (
	"GoAuth/internal/adapters/http/dto"
	"GoAuth/internal/adapters/http/validation"
	"GoAuth/internal/ports/inbounds"
	"net/http"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
)

type SchemaHandler struct {
	schemas inbounds.SchemaService
}

func NewSchemaHandler(uc inbounds.SchemaService) *SchemaHandler {
	return &SchemaHandler{schemas: uc}
}

// Draft godoc
// @Summary Draft a new schema
// @Description Creates a new schema draft for a project.
// @Tags schemas
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Param project_id path string true "Project ID"
// @Param schemaInfo body dto.DraftSchemaRequest true "Draft Schema Request"
// @Success 201 {object} dto.SchemaResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects/{project_id}/schemas [post]
func (handler *SchemaHandler) Draft(w http.ResponseWriter, r *http.Request) {
	projectID, rs := getUUID(r, "project_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	var req dto.DraftSchemaRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		resp.FromError(err).Send(w)
		return
	}

	in := inbounds.SchemaServiceInput{
		SchemaType: req.SchemaType,
		Title:      req.Title,
		FlowID:     req.FlowID,
		ProjectID:  projectID,
	}

	ctx := r.Context()
	res, err := handler.schemas.Draft(ctx, in)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.Created("drafted schema").
		WithData(dto.SchemaOutputToResponse(res)).
		Send(w)
}

// Publish godoc
// @Summary Publish a schema
// @Description Publishes a schema for a project.
// @Tags schemas
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Param project_id path string true "Project ID"
// @Param schema_id path string true "Schema ID"
// @Success 200 {object} object
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects/{project_id}/schemas/{schema_id}/publish [post]
func (handler *SchemaHandler) Publish(w http.ResponseWriter, r *http.Request) {
	projectID, rs := getUUID(r, "project_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	schemaID, rs := getUUID(r, "schema_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	in := inbounds.SchemaServiceInput{
		ProjectID: projectID,
		SchemaID:  schemaID,
	}

	ctx := r.Context()
	if err := handler.schemas.Publish(ctx, in); err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK("published schema").Send(w)
}

// GetByID godoc
// @Summary Get a schema by ID
// @Description Retrieves a schema by its ID.
// @Tags schemas
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Param project_id path string true "Project ID"
// @Param schema_id path string true "Schema ID"
// @Success 200 {object} dto.SchemaResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects/{project_id}/schemas/{schema_id} [get]
func (handler *SchemaHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	projectID, rs := getUUID(r, "project_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	schemaID, rs := getUUID(r, "schema_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	in := inbounds.SchemaServiceInput{
		ProjectID: projectID,
		SchemaID:  schemaID,
	}

	ctx := r.Context()
	found, err := handler.schemas.GetByID(ctx, in)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK().
		WithData(dto.SchemaOutputToResponse(found)).
		Send(w)
}

// GetVerbose godoc
// @Summary Get a verbose schema by ID
// @Description Retrieves a verbose schema by its ID.
// @Tags schemas
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Param project_id path string true "Project ID"
// @Param schema_id path string true "Schema ID"
// @Success 200 {object} dto.VerboseSchemaResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects/{project_id}/schemas/{schema_id}/verbose [get]
func (handler *SchemaHandler) GetVerbose(w http.ResponseWriter, r *http.Request) {
	projectID, rs := getUUID(r, "project_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	schemaID, rs := getUUID(r, "schema_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	in := inbounds.SchemaServiceInput{
		ProjectID: projectID,
		SchemaID:  schemaID,
	}

	ctx := r.Context()
	res, err := handler.schemas.GetVerbose(ctx, in)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK().
		WithData(dto.VerboseSchemaOutputToResponse(res)).
		Send(w)
}
