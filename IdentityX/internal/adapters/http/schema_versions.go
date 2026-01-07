package http

import (
	"GoAuth/internal/adapters/http/dto"
	"GoAuth/internal/ports/inbounds"
	"net/http"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/go-chi/chi/v5"
)

type SchemaVersionHandler struct {
	versions inbounds.SchemaVersionService
}

func NewSchemaVersionHandler(uc inbounds.SchemaVersionService) *SchemaVersionHandler {
	return &SchemaVersionHandler{versions: uc}
}

func (handler *SchemaVersionHandler) Draft(w http.ResponseWriter, r *http.Request) {
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

	in := inbounds.DraftSchemaVersionInput{
		ProjectID: projectID,
		SchemaID:  schemaID,
	}

	ctx := r.Context()
	res, err := handler.versions.Draft(ctx, in)
	if err != nil {
		ErrToResp(err).Send(w)
		return
	}

	resp.Created("drafted schema version").
		WithData(dto.SchemaVersionOutputToResponse(res)).
		Send(w)
}

func (handler *SchemaVersionHandler) Publish(w http.ResponseWriter, r *http.Request) {
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

	in := inbounds.PublishSchemaVersionInput{
		ProjectID: projectID,
		SchemaID:  schemaID,
	}

	ctx := r.Context()
	err := handler.versions.Publish(ctx, in)
	if err != nil {
		ErrToResp(err).Send(w)
		return
	}

	resp.OK("published schema version").Send(w)
}
