package handlers

import (
	"GoAuth/internal/adapters/http/dto"
	"GoAuth/internal/adapters/http/validation"
	"GoAuth/internal/ports/inbounds"
	"net/http"
	"strconv"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/go-chi/chi/v5"
)

type SchemaFieldsHandler struct {
	fields inbounds.SchemaFieldsService
}

func NewSchemaFieldsHandler(uc inbounds.SchemaFieldsService) *SchemaFieldsHandler {
	return &SchemaFieldsHandler{fields: uc}
}

// Create godoc
// @Summary Create fields for a schema version
// @Description Creates fields for a specific version of a schema.
// @Tags schema-fields
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Param project_id path string true "Project ID"
// @Param schema_id path string true "Schema ID"
// @Param version path int true "Schema Version Number"
// @Param fieldInfo body dto.CreateFieldRequest true "Field creation information"
// @Success 201 {array} dto.FieldResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects/{project_id}/schemas/{schema_id}/v{version} [post]
func (handler *SchemaFieldsHandler) Create(w http.ResponseWriter, r *http.Request) {
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

	version := chi.URLParam(r, "version")
	if version == "" {
		resp.BadRequest("missing version parameter").Send(w)
		return
	}

	versionNumber, err := strconv.Atoi(version)
	if err != nil {
		resp.BadRequest("invalid version parameter").AddTrace(err).Send(w)
		return
	}

	if versionNumber <= 0 {
		resp.BadRequest("version must be >= 1").Send(w)
		return
	}

	var req dto.CreateFieldRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		resp.FromError(err).Send(w)
		return
	}

	in := inbounds.SchemaFieldInput{
		ProjectID:     projectID,
		SchemaID:      schemaID,
		VersionNumber: versionNumber,
		Fields:        dto.FieldParamSliceToInputFieldSlice(req.Fields),
	}

	ctx := r.Context()
	res, err := handler.fields.Create(ctx, in)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.Created("created fields").
		WithData(dto.OutputFieldSliceToFieldResponseSlice(res)).
		Send(w)
}
