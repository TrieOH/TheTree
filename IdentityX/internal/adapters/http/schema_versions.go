package http

import (
	"GoAuth/internal/adapters/http/dto"
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
