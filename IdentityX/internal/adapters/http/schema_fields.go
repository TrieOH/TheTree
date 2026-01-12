package http

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

func (handler *SchemaFieldsHandler) Create(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "project_id")
	if projectID == "" {
		resp.BadRequest("missing project id parameter").Send(w)
		return
	}

	schemaID := chi.URLParam(r, "schema_id")
	if schemaID == "" {
		resp.BadRequest("missing schema id parameter").Send(w)
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
