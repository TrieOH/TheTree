package handlers

import (
	"GoAuth/internal/adapters/http/dto"
	"GoAuth/internal/adapters/http/validation"
	"GoAuth/internal/ports/inbounds"
	"net/http"
	"strconv"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
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

// GetIDsFromProjectID godoc
// @Summary Get schema IDs from a project
// @Description Retrieves all schema IDs for a given project.
// @Tags schemas
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Param project_id path string true "Project ID"
// @Success 200 {array} string
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects/{project_id}/schemas/ids [get]
func (handler *SchemaHandler) GetIDsFromProjectID(w http.ResponseWriter, r *http.Request) {
	projectID, rs := getUUID(r, "project_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	ctx := r.Context()
	IDs, err := handler.schemas.GetIDsFromProjectID(ctx, projectID)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK().WithData(IDs).Send(w)
}

// List godoc
// @Summary List schemas
// @Description Retrieves all schemas for a given project.
// @Tags schemas
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Param project_id path string true "Project ID"
// @Success 200 {array} dto.SchemaResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects/{project_id}/schemas [get]
func (handler *SchemaHandler) List(w http.ResponseWriter, r *http.Request) {
	projectID, rs := getUUID(r, "project_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	ctx := r.Context()
	schemas, err := handler.schemas.List(ctx, projectID)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK().WithData(dto.SchemaOutputSliceToResponse(schemas)).Send(w)
}

// GetLatestFormByID godoc
// @Summary Get latest form by schema ID
// @Description Retrieves the current published version of a specific schema (UUID)
// @Tags schemas
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Param project_id path string true "Project ID"
// @Param schema_id path string true "Schema UUID"
// @Success 200 {object} dto.FormResponse
// @Router /projects/{project_id}/schemas/{schema_id}/latest [get]
func (handler *SchemaHandler) GetLatestFormByID(w http.ResponseWriter, r *http.Request) {
	projectID, rs := getUUID(r, "project_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	schemaID, rs := getUUID(r, "schema_id")
	if rs != nil {
		resp.BadRequest("invalid schema_id format").Send(w)
		return
	}

	in := inbounds.SchemaServiceInput{
		ProjectID: projectID,
		SchemaID:  schemaID,
	}

	ctx := r.Context()
	form, err := handler.schemas.GetLatestForm(ctx, in)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK().WithData(dto.FormOutputToFormToResponse(form)).Send(w)
}

// GetLatestFormByFlow godoc
// @Summary Get latest form by flow ID and type
// @Description Retrieves the current published version by flow identifier
// @Tags schemas
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Param project_id path string true "Project ID"
// @Param flow_id query string true "Flow ID (e.g., estudante, professor)"
// @Param schema_type query string false "Schema type" Enums(core, context, sub-context) default(context)
// @Success 200 {object} dto.FormResponse
// @Router /projects/{project_id}/schemas/lookup/latest [get]
func (handler *SchemaHandler) GetLatestFormByFlow(w http.ResponseWriter, r *http.Request) {
	projectID, rs := getUUID(r, "project_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	flowID := r.URL.Query().Get("flow_id")
	if flowID == "" {
		resp.BadRequest("flow_id query parameter required").Send(w)
		return
	}

	schemaType := r.URL.Query().Get("schema_type")
	if schemaType == "" {
		schemaType = "context"
	}

	in := inbounds.SchemaServiceInput{
		ProjectID:  projectID,
		SchemaID:   uuid.Nil, // nil signals lookup by flow+type
		FlowID:     flowID,
		SchemaType: schemaType,
	}

	ctx := r.Context()
	form, err := handler.schemas.GetLatestForm(ctx, in)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK().WithData(dto.FormOutputToFormToResponse(form)).Send(w)
}

// GetSpecificFormByID handles versioned form by UUID
func (handler *SchemaHandler) GetSpecificFormByID(w http.ResponseWriter, r *http.Request) {
	projectID, rs := getUUID(r, "project_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	schemaID, rs := getUUID(r, "schema_id")
	if rs != nil {
		resp.BadRequest("invalid schema_id format").Send(w)
		return
	}

	versionStr := chi.URLParam(r, "version")
	versionNum, err := strconv.Atoi(versionStr)
	if err != nil || versionNum < 1 {
		resp.BadRequest("invalid version number").Send(w)
		return
	}

	in := inbounds.SchemaServiceInput{
		ProjectID: projectID,
		SchemaID:  schemaID,
	}

	ctx := r.Context()
	form, err := handler.schemas.GetFormByVersion(ctx, in, versionNum)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK().WithData(dto.FormOutputToFormToResponse(form)).Send(w)
}

// GetSpecificFormByFlow handles versioned form by flow lookup
func (handler *SchemaHandler) GetSpecificFormByFlow(w http.ResponseWriter, r *http.Request) {
	projectID, rs := getUUID(r, "project_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	flowID := r.URL.Query().Get("flow_id")
	if flowID == "" {
		resp.BadRequest("flow_id query parameter required").Send(w)
		return
	}

	schemaType := r.URL.Query().Get("schema_type")
	if schemaType == "" {
		schemaType = "context"
	}

	versionStr := chi.URLParam(r, "version")
	versionNum, err := strconv.Atoi(versionStr)
	if err != nil || versionNum < 1 {
		resp.BadRequest("invalid version number").Send(w)
		return
	}

	in := inbounds.SchemaServiceInput{
		ProjectID:  projectID,
		SchemaID:   uuid.Nil,
		FlowID:     flowID,
		SchemaType: schemaType,
	}

	ctx := r.Context()
	form, err := handler.schemas.GetFormByVersion(ctx, in, versionNum)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK().WithData(dto.FormOutputToFormToResponse(form)).Send(w)
}
