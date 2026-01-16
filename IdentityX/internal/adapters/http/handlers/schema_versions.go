package handlers

import (
	"GoAuth/internal/adapters/http/dto"
	"GoAuth/internal/adapters/http/validation"
	"GoAuth/internal/ports/inbounds"
	"net/http"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
)

type SchemaVersionHandler struct {
	versions inbounds.SchemaVersionService
}

func NewSchemaVersionHandler(uc inbounds.SchemaVersionService) *SchemaVersionHandler {
	return &SchemaVersionHandler{versions: uc}
}

// Draft godoc
// @Summary Draft a new schema version
// @Description Creates a new version draft for a schema.
// @Tags schema-versions
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Param project_id path string true "Project ID"
// @Param schema_id path string true "Schema ID"
// @Success 201 {object} dto.SchemaVersionResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects/{project_id}/schemas/{schema_id}/versions/draft [post]
func (handler *SchemaVersionHandler) Draft(w http.ResponseWriter, r *http.Request) {
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

	in := inbounds.SchemaVersionServiceInput{
		ProjectID: projectID,
		SchemaID:  schemaID,
	}

	ctx := r.Context()
	res, err := handler.versions.Draft(ctx, in)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.Created("drafted schema version").
		WithData(dto.SchemaVersionOutputToResponse(res)).
		Send(w)
}

// Publish godoc
// @Summary Publish a schema version
// @Description Publishes a version of a schema.
// @Tags schema-versions
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Param project_id path string true "Project ID"
// @Param schema_id path string true "Schema ID"
// @Success 200 {object} object
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects/{project_id}/schemas/{schema_id}/versions/publish [post]
func (handler *SchemaVersionHandler) Publish(w http.ResponseWriter, r *http.Request) {
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

	in := inbounds.SchemaVersionServiceInput{
		ProjectID: projectID,
		SchemaID:  schemaID,
	}

	ctx := r.Context()
	err := handler.versions.Publish(ctx, in)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK("published schema version").Send(w)
}

// GetCurrent godoc
// @Summary Get the current schema version
// @Description Retrieves the current version of a schema.
// @Tags schema-versions
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Param project_id path string true "Project ID"
// @Param schema_id path string true "Schema ID"
// @Success 200 {object} dto.SchemaVersionResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects/{project_id}/schemas/{schema_id}/versions/current [get]
func (handler *SchemaVersionHandler) GetCurrent(w http.ResponseWriter, r *http.Request) {
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

	in := inbounds.SchemaVersionServiceInput{
		SchemaID:  schemaID,
		ProjectID: projectID,
	}

	ctx := r.Context()
	current, err := handler.versions.GetCurrent(ctx, in)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK().WithData(dto.SchemaVersionOutputToResponse(current)).Send(w)
}

// GetLatest godoc
// @Summary Get the latest schema version
// @Description Retrieves the latest version of a schema.
// @Tags schema-versions
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Param project_id path string true "Project ID"
// @Param schema_id path string true "Schema ID"
// @Success 200 {object} dto.SchemaVersionResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects/{project_id}/schemas/{schema_id}/versions/latest [get]
func (handler *SchemaVersionHandler) GetLatest(w http.ResponseWriter, r *http.Request) {
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

	in := inbounds.SchemaVersionServiceInput{
		SchemaID:  schemaID,
		ProjectID: projectID,
	}

	ctx := r.Context()
	latest, err := handler.versions.GetLatest(ctx, in)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK().WithData(dto.SchemaVersionOutputToResponse(latest)).Send(w)
}

// GetVerbose godoc
// @Summary Get a verbose schema version
// @Description Retrieves a verbose version of a schema.
// @Description Prefers you send version ID when available, version number is mandatory
// @Tags schema-versions
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Param project_id path string true "Project ID"
// @Param schema_id path string true "Schema ID"
// @Param version path int true "Version Number"
// @Param versionInfo body dto.GetVersionVerboseRequest true "VersionID if available"
// @Success 200 {object} dto.VersionVerboseResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects/{project_id}/schemas/{schema_id}/versions/v{version} [get]
func (handler *SchemaVersionHandler) GetVerbose(w http.ResponseWriter, r *http.Request) {
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

	versionNumber, rs := getNumber(r, "version")
	if rs != nil {
		rs.Send(w)
		return
	}

	var req dto.GetVersionVerboseRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		resp.FromError(err).Send(w)
		return
	}

	in := inbounds.SchemaVersionServiceInput{
		ProjectID:     projectID,
		SchemaID:      schemaID,
		VersionID:     req.VersionID,
		VersionNumber: versionNumber,
	}

	ctx := r.Context()
	version, err := handler.versions.GetVerbose(ctx, in)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK().WithData(dto.VerboseVersionOutputToResponse(version)).Send(w)
}
