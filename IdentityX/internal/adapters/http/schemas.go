package http

import (
	"GoAuth/internal/adapters/http/dto"
	"GoAuth/internal/ports/inbounds"
	"net/http"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/MintzyG/FastUtilitiesNet/validation"
	"github.com/go-chi/chi/v5"
)

type SchemaHandler struct {
	schemas inbounds.SchemaService
}

func NewSchemaHandler(uc inbounds.SchemaService) *SchemaHandler {
	return &SchemaHandler{schemas: uc}
}

func (handler *SchemaHandler) Draft(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "project_id")
	if projectID == "" {
		resp.BadRequest("missing project id parameter").Send(w)
		return
	}

	var req dto.DraftSchemaRequest
	if rs := validation.ValidateInto(r, &req); rs != nil {
		rs.Send(w)
		return
	}

	in := inbounds.DraftSchemaInput{
		SchemaType: req.SchemaType,
		Title:      req.Title,
		FlowID:     req.FlowID,
		ProjectID:  projectID,
	}

	ctx := r.Context()
	res, err := handler.schemas.Draft(ctx, in)
	if err != nil {
		ErrToResp(err).Send(w)
		return
	}

	resp.Created("drafted schema").
		WithData(dto.SchemaOutputToResponse(res)).
		Send(w)
}

func (handler *SchemaHandler) Publish(w http.ResponseWriter, r *http.Request) {
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

	in := inbounds.PublishSchemaInput{
		ProjectID: projectID,
		SchemaID:  schemaID,
	}

	ctx := r.Context()
	if err := handler.schemas.Publish(ctx, in); err != nil {
		ErrToResp(err).Send(w)
		return
	}

	resp.OK("published schema").Send(w)
}

func (handler *SchemaHandler) GetByID(w http.ResponseWriter, r *http.Request) {
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

	in := inbounds.GetSchemaByIDInput{
		ProjectID: projectID,
		SchemaID:  schemaID,
	}

	ctx := r.Context()
	found, err := handler.schemas.GetByID(ctx, in)
	if err != nil {
		ErrToResp(err).Send(w)
		return
	}

	resp.OK().
		WithData(dto.SchemaOutputToResponse(found)).
		Send(w)
}
