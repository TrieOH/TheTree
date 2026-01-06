package http

import (
	"GoAuth/internal/adapters/http/dto"
	"GoAuth/internal/ports/inbounds"
	"net/http"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/MintzyG/FastUtilitiesNet/validation"
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

	var req dto.DraftSchemaVersionRequest
	if rs := validation.ValidateInto(r, &req); rs != nil {
		rs.Send(w)
		return
	}

	in := inbounds.DraftSchemaVersionInput{
		ProjectID: projectID,
		SchemaID:  req.SchemaID,
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
