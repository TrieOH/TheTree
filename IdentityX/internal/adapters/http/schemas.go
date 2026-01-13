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
